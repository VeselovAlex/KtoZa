package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/VeselovAlex/KtoZa/model"
)

type AnswerController struct{}

func (ctrl *AnswerController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metod not allowed. Use method POST instead", http.StatusMethodNotAllowed)
		return
	}
	var answers model.AnswerSet
	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		log.Println("ANSWER", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	StatisticsStorage.ApplyAnswerSet(answers)
}
