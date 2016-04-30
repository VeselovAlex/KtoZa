package main

import (
	"encoding/json"
	"fmt"
<<<<<<< HEAD
	"log"
	"net/http"
)

=======
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

>>>>>>> ef16452... poll remodelled, server updated
// PollController обрабатывает запросы, связанные с данными опроса
type PollController struct {
}

<<<<<<< HEAD
func (ctrl *PollController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctrl.handleGetPoll(w, r)
	case http.MethodPut:
		ctrl.handleCreateOrUpdatePoll(w, r)
=======
func (p *PollController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.handleGetPoll(w, r)
>>>>>>> ef16452... poll remodelled, server updated
	default:
		errMsg := fmt.Sprint("Method &s in unsupported", r.Method)
		http.Error(w, errMsg, http.StatusMethodNotAllowed)
	}
}

<<<<<<< HEAD
func (ctrl *PollController) handleGetPoll(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(App.PollStorage.Get())
=======
func (p *PollController) handleGetPoll(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(dumb)
>>>>>>> ef16452... poll remodelled, server updated
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
<<<<<<< HEAD

func (ctrl *PollController) handleCreateOrUpdatePoll(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateOrUpdate call")
	poll := App.PollStorage.Get()
	var err error
	// poll := &Poll{}
	// err := App.RequestDecoder.FromRequest(r, poll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	App.PollStorage.CreateOrUpdate(poll)
	Respondents.NotifyAll(About.UpdatedPoll(poll))
}
=======
>>>>>>> ef16452... poll remodelled, server updated
