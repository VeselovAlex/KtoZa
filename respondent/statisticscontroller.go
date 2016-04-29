package main

import (
	"encoding/json"
	"log"
	"net/http"
<<<<<<< HEAD
)

// StatisticsController обрабатывает запросы, связанные с данными статистики для текущего опроса
=======

	"github.com/VeselovAlex/KtoZa/model"
)

>>>>>>> 6e99eb4... session controller, app configuration added
type StatisticsController struct{}

func (ctrl *StatisticsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metod not allowed. Use method GET instead", http.StatusMethodNotAllowed)
		return
	}
<<<<<<< HEAD
	err := json.NewEncoder(w).Encode(App.StatisticsStorage.Get())
=======
	dumbStat := model.CreateStatisticsFor(App.PollStorage.Get())
	err := json.NewEncoder(w).Encode(dumbStat)
>>>>>>> 6e99eb4... session controller, app configuration added
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error serving poll:", err)
	}
}
