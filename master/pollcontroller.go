package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/VeselovAlex/KtoZa/model"
)

var dumbPoll = &model.Poll{
	Title:   "Dummy poll",
	Caption: "Простой опрос для тестирования",
	Events: model.EventTimings{
		RegistrationAt: time.Now().Add(5 * time.Second),
		StartAt:        time.Now().Add(35 * time.Second),
		EndAt:          time.Now().Add(95 * time.Second),
	},
	Questions: []model.Question{
		model.Question{
			Text:    "Вопрос 1",
			Type:    "single-option",
			Options: []string{"Yes", "No"},
		},
		model.Question{
			Text:    "Вопрос 2",
			Type:    "multi-option",
			Options: []string{"One", "Two", "Three"},
		},
	},
}

// PollController обрабатывает запросы, связанные с данными опроса
type PollController struct {
}

func (ctrl *PollController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctrl.handleGetPoll(w, r)
	case http.MethodPut:
		ctrl.handleCreateOrUpdatePoll(w, r)
	default:
		errMsg := fmt.Sprint("Method &s in unsupported", r.Method)
		http.Error(w, errMsg, http.StatusMethodNotAllowed)
	}
}

func (ctrl *PollController) handleGetPoll(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(dumbPoll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ctrl *PollController) handleCreateOrUpdatePoll(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateOrUpdate call")
}
