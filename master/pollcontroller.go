package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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
	err := json.NewEncoder(w).Encode(App.PollStorage.Get())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
