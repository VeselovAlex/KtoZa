package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/VeselovAlex/KtoZa/model"
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

func NewStatisticsController(listeners ...StatisticsListener) *StatisticsController {
	return &StatisticsController{
		listeners: listeners,
	}
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

func (ctrl *StatisticsController) OnPollUpdate(poll *model.Poll) {
	ctrl.lock.Lock()
	defer ctrl.lock.Unlock()

	ctrl.stat = model.CreateStatisticsFor(poll)
	ctrl.notifyListeners(ctrl.stat)
}
