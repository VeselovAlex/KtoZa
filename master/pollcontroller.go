package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VeselovAlex/KtoZa/model"
)

var dumb = &model.Poll{
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

func (p *PollController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.handleGetPoll(w, r)
	default:
		errMsg := fmt.Sprint("Method &s in unsupported", r.Method)
		http.Error(w, errMsg, http.StatusMethodNotAllowed)
	}
}

func (p *PollController) handleGetPoll(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(dumb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}