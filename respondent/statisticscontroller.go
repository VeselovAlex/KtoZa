package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// StatisticsController обрабатывает запросы, связанные с данными статистики для текущего опроса
type StatisticsController struct{}

func (ctrl *StatisticsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metod not allowed. Use method GET instead", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewEncoder(w).Encode(App.StatisticsStorage.Get())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error serving poll:", err)
	}
}
