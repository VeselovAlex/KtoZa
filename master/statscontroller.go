package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VeselovAlex/KtoZa/model"
)

// StatisticsController обрабатывает запросы, связанные с данными статистики
type StatisticsController struct{}

func (ctrl *StatisticsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctrl.handleGetStats(w, r)
	default:
		errMsg := fmt.Sprint("Method &s in unsupported", r.Method)
		http.Error(w, errMsg, http.StatusMethodNotAllowed)
	}
}

// For debug purposes
func (ctrl *StatisticsController) handleGetStats(w http.ResponseWriter, r *http.Request) {
	dumbStats := model.CreateStatisticsFor(App.PollStorage.Get())
	err := json.NewEncoder(w).Encode(dumbStats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
