package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VeselovAlex/KtoZa/model"
)

// StatisticsController обрабатывает запросы, связанные с данными статистики
type StatisticsController struct{}

func (ctrl *StatisticsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctrl.handleGetStats(w, r)
	case http.MethodDelete:
		ctrl.handleDeleteStats(w, r)
	default:
		errMsg := fmt.Sprint("Method &s in unsupported", r.Method)
		http.Error(w, errMsg, http.StatusMethodNotAllowed)
		log.Printf("STATISTICS [%s] :: %s\n", r.Method, errMsg)
	}
}

func (ctrl *StatisticsController) handleGetStats(w http.ResponseWriter, r *http.Request) {
	stat := App.StatisticsStorage.Get()
	if stat == nil {
		// Если статистика не создана, то пытаемся создать
		poll := App.PollStorage.Get()
		if poll != nil {
			// Если задан опрос, то создаем статистику
			stat = model.CreateStatisticsFor(poll)
			App.StatisticsStorage.CreateOrUpdate(stat)
		}
	}

	// Кодируем статистику в JSON и отправляем
	err := json.NewEncoder(w).Encode(stat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("STATISTICS [GET] :: Error:", err)
	}
}

func (ctrl *StatisticsController) handleDeleteStats(w http.ResponseWriter, r *http.Request) {
	stat := App.StatisticsStorage.Delete()
	if stat != nil {
		// Статистика не была удалена ранее
		App.PubSub.NotifyAll(About.UpdatedStatistics(nil))
		log.Println("STATISTICS [DELETE] :: Statistics delete")
	}
}
