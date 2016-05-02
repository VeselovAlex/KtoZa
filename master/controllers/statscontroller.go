package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/VeselovAlex/KtoZa/model"

	common "github.com/VeselovAlex/KtoZa"
)

type StatisticsListener interface {
	OnStatisticsUpdate(*model.Statistics)
}

// StatisticsController обрабатывает запросы, связанные с данными статистики
type StatisticsController struct {
	listeners []StatisticsListener

	lock sync.RWMutex
	stat *model.Statistics
}

// NewStatisticsController создает новый экземпляр контроллера статистики
func NewStatisticsController(listeners ...StatisticsListener) *StatisticsController {
	listeners = append(listeners, Storage)
	ctrl := &StatisticsController{
		listeners: listeners,
	}
	go ctrl.listenForAnswerCaches()
	return ctrl
}

func (ctrl *StatisticsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctrl.handleGetStats(w, r)
	case http.MethodDelete:
		ctrl.handleDeleteStats(w, r)
	default:
		errMsg := fmt.Sprint("Method &s in unsupported", r.Method)
		http.Error(w, errMsg, http.StatusMethodNotAllowed)
		log.Printf("STAT CONTROLLER :: [%s] %s\n", r.Method, errMsg)
	}
}

func (ctrl *StatisticsController) handleGetStats(w http.ResponseWriter, r *http.Request) {
	ctrl.lock.RLock()
	defer ctrl.lock.RUnlock()

	// Кодируем статистику в JSON и отправляем
	err := json.NewEncoder(w).Encode(ctrl.stat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("STATISTICS [GET] :: Error:", err)
	}
}

func (ctrl *StatisticsController) handleDeleteStats(w http.ResponseWriter, r *http.Request) {
	ctrl.lock.Lock()
	defer ctrl.lock.Unlock()
	stat := ctrl.stat
	if stat != nil {
		// Статистика не была удалена ранее
		stat = nil
		ctrl.notifyListeners(nil)
		log.Println("STATISTICS [DELETE] :: Statistics delete")
	}
}

func (ctrl *StatisticsController) notifyListeners(stat *model.Statistics) {
	for _, listener := range ctrl.listeners {
		listener.OnStatisticsUpdate(stat)
	}
}

// OnPollUpdate -- действие контроллера статистики при изменении опроса
func (ctrl *StatisticsController) OnPollUpdate(poll *model.Poll) {
	ctrl.lock.Lock()
	defer ctrl.lock.Unlock()
	if poll == nil {
		ctrl.stat = nil
	} else if ctrl.stat == nil {
		// Первоначальная загрузка
		var err error
		stat, err := Storage.ReadStatistics()
		newStat := model.CreateStatisticsFor(poll)
		if err != nil || stat == nil || !newStat.IsJoinableWith(stat) {
			//Нет статистики или статистика не соответствует опросу
			ctrl.stat = newStat
		} else {
			ctrl.stat = stat
		}
	} else {
		ctrl.stat = model.CreateStatisticsFor(poll)
	}
	ctrl.notifyListeners(ctrl.stat)
}

func (ctrl *StatisticsController) listenForAnswerCaches() {
	apply := func(cache *model.Statistics) {
		ctrl.lock.Lock()
		defer ctrl.lock.Unlock()
		if ctrl.stat == nil {
			// Пропускаем
			return
		}
		if ok := ctrl.stat.JoinWith(cache); ok {
			ctrl.notifyListeners(ctrl.stat)
		}
	}
	for {
		msg, ok := Respondents.Await().(common.EventRawMessage)
		if !ok || msg.Event != common.EventNewAnswerCache {
			log.Println("STATISTICS CTRL :: Bad message got from respondent")
			if !ok {
				log.Println("STATISTICS CTRL :: Bad message casting")
			}
			continue
		}
		cache := &model.Statistics{}
		err := json.Unmarshal(msg.Data, cache)
		if err != nil {
			log.Println("STATISTICS CTRL :: Bad message got from respondent:", err)
			continue
		}
		apply(cache)
		log.Println("STATISTICS CTRL :: Applied answer")
	}
}
