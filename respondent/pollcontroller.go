package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type PollController struct{}

func (ctrl *PollController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metod not allowed. Use method GET instead", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewEncoder(w).Encode(App.PollStorage.Get())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error serving poll:", err)
	}
}
