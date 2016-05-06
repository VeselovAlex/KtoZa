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

	cacheLock  sync.RWMutex
	needUpdate bool
	cache      *model.Statistics

	snapshotLock sync.RWMutex
	snapshot     *model.Statistics
}

func (ctrl *StatisticsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metod not allowed. Use method GET instead", http.StatusMethodNotAllowed)
		return
	}

	encoder := json.NewEncoder(w)
	err := func() error {
		ctrl.snapshotLock.RLock()
		defer ctrl.snapshotLock.RUnlock()
		return encoder.Encode(ctrl.snapshot)
	}()
	if err != nil {
		log.Println("STAT CONTROLLER :: Error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Println("STAT CONTROLLER :: Success")
	}
}

func NewStatisticsController(listeners ...StatisticsListener) *StatisticsController {
	ctrl := &StatisticsController{
		listeners: listeners,
	}
	go ctrl.listenToMaster()
	go ctrl.doSync()
	return ctrl
}

// OnPollUpdate обновляет состояние контроллера при изменении опроса
func (ctrl *StatisticsController) OnPollUpdate(poll *model.Poll) {
	ctrl.cacheLock.Lock()
	defer ctrl.cacheLock.Unlock()
	ctrl.snapshotLock.Lock()
	defer ctrl.snapshotLock.Unlock()
	ctrl.cache = model.CreateStatisticsFor(poll)
	stats, err := MasterServer.GetStatistics()
	if err != nil {
		log.Fatalln("STAT CONTROLLER :: Bad poll update:", err)
	}
	ctrl.snapshot = stats
	ctrl.onSnapshotUpdated()
}

// OnNewAnswerSet обновляет состояние контроллера при получении ответа
func (ctrl *StatisticsController) OnNewAnswerSet(ans model.AnswerSet) {
	ctrl.cacheLock.Lock()
	defer ctrl.cacheLock.Unlock()
	applied := ctrl.cache.ApplyAnswerSet(ans)
	if applied {
		ctrl.needUpdate = true
	}
}

func (ctrl *StatisticsController) doSync() {
	moveIfNeeded := func() (*model.Statistics, bool) {
		var cpy *model.Statistics
		ctrl.cacheLock.RLock()
		defer ctrl.cacheLock.RUnlock()
		need := ctrl.needUpdate
		if need {
			// Копируем текущую версию кэша
			cpy = &model.Statistics{}
			ctrl.cache.CopyTo(cpy)

			// Обнуляем кеш
			ctrl.cache.LastUpdate = time.Now()
			ctrl.cache.RespondentsCount = 0
			questions := ctrl.cache.Questions
			for i := range questions {
				questions[i].AnswersCount = 0
				opts := questions[i].Options
				for j := range opts {
					opts[j].Count = 0
				}
			}

			ctrl.needUpdate = false
		}
		return cpy, need
	}
	for {
		cpy, upd := moveIfNeeded()
		if upd {
			err := MasterServer.SendAnswerCache(cpy)
			if err != nil {
				// Если не удалось обновиться, нужно объединить текущий кэш с сохраненным
				func() {
					ctrl.cacheLock.Lock()
					defer ctrl.cacheLock.Unlock()
					ctrl.cache.JoinWith(cpy)
					ctrl.needUpdate = true
				}()
				// Ждем секунду перед новой попыткой
				time.Sleep(time.Second)
			}
		}
	}
}

func (ctrl *StatisticsController) listenToMaster() {
	updateWith := func(snapshot *model.Statistics) {
		ctrl.snapshotLock.Lock()
		defer ctrl.snapshotLock.Unlock()
		if ctrl.snapshot.LastUpdate.Before(snapshot.LastUpdate) {
			ctrl.snapshot = snapshot
			log.Println("STAT CONTROLLER :: Fetched snapshot from master")
			ctrl.onSnapshotUpdated()
		}
	}
	for {
		snapshot := MasterServer.AwaitStatisticsUpdate()
		updateWith(snapshot)
	}
}

func (ctrl *StatisticsController) onSnapshotUpdated() {
	ctrl.notifyListeners(ctrl.snapshot)
}

func (ctrl *StatisticsController) notifyListeners(stat *model.Statistics) {
	for _, listener := range ctrl.listeners {
		listener.OnStatisticsUpdate(stat)
	}
}
