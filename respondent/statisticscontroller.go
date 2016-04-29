package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/VeselovAlex/KtoZa/model"
)

type StatisticsController struct{}

func (ctrl *StatisticsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metod not allowed. Use method GET instead", http.StatusMethodNotAllowed)
		return
	}
	dumbStat := model.CreateStatisticsFor(App.PollStorage.Get())
	err := json.NewEncoder(w).Encode(dumbStat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error serving poll:", err)
	}
}
