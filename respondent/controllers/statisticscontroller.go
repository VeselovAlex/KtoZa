package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/VeselovAlex/KtoZa/model"
)

type StatisticsListener interface {
	OnStatisticsUpdate(*model.Statistics)
}

// StatisticsController обрабатывает запросы, связанные с данными статистики для текущего опроса
type StatisticsController struct {
	listeners []StatisticsListener

	lock sync.RWMutex
	stat *model.Statistics
}

func (ctrl *StatisticsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metod not allowed. Use method GET instead", http.StatusMethodNotAllowed)
		return
	}

	encoder := json.NewEncoder(w)
	err := func() error {
		ctrl.lock.RLock()
		defer ctrl.lock.RUnlock()
		return encoder.Encode(ctrl.stat)
	}()
	if err != nil {
		log.Println("STAT CONTROLLER :: Error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Println("STAT CONTROLLER :: Success")
	}
}

func NewTestStatCtrl(listeners ...StatisticsListener) *StatisticsController {
	return &StatisticsController{
		listeners: listeners,
	}
}

func (ctrl *StatisticsController) OnPollUpdate(poll *model.Poll) {
	ctrl.lock.Lock()
	defer ctrl.lock.Unlock()
	ctrl.stat = &model.Statistics{
		LastUpdate: time.Now(),
	}
	ctrl.notifyListeners(ctrl.stat)
}

func (ctrl *StatisticsController) OnNewAnswerSet(ans model.AnswerSet) {
	ctrl.lock.Lock()
	defer ctrl.lock.Unlock()
	ctrl.stat = &model.Statistics{
		LastUpdate: time.Now(),
	}
	ctrl.notifyListeners(ctrl.stat)
}

func (ctrl *StatisticsController) notifyListeners(stat *model.Statistics) {
	for _, listener := range ctrl.listeners {
		listener.OnStatisticsUpdate(stat)
	}
}
