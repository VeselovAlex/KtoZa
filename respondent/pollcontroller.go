package main

import (
	"encoding/json"
	"log"
	"net/http"
)

<<<<<<< HEAD
<<<<<<< HEAD
// PollController обрабатывает запросы, связанные с данными опроса
=======
>>>>>>> 6e99eb4... session controller, app configuration added
=======
// PollController обрабатывает запросы, связанные с данными опроса
>>>>>>> 2c58c96... Added Statistics storage
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
