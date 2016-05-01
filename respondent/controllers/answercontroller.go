package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/VeselovAlex/KtoZa/model"
)

type AnswerListener interface {
	OnNewAnswerSet(model.AnswerSet)
}

type AnswerController struct {
	listeners []AnswerListener
}

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
	ctrl.notifyListeners(answers)
}

func (ctrl *AnswerController) notifyListeners(ans model.AnswerSet) {
	for _, listener := range ctrl.listeners {
		listener.OnNewAnswerSet(ans)
	}
}

func NewTestAnswerCtrl(listeners ...AnswerListener) *AnswerController {
	ctrl := &AnswerController{
		listeners: listeners,
	}
	return ctrl
}
