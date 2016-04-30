package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VeselovAlex/KtoZa/model"
)

// PollController обрабатывает запросы, связанные с данными опроса
type PollController struct{}

func (ctrl *PollController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctrl.handleGetPoll(w, r)
	case http.MethodPut:
		ctrl.handleCreateOrUpdatePoll(w, r)
	case http.MethodDelete:
		ctrl.handleDeletePoll(w, r)
	default:
		errMsg := fmt.Sprint("Method &s in unsupported", r.Method)
		http.Error(w, errMsg, http.StatusMethodNotAllowed)
		log.Printf("POLL [%s] :: %s\n", r.Method, errMsg)
	}
}

// GET
func (ctrl *PollController) handleGetPoll(w http.ResponseWriter, r *http.Request) {
	poll := App.PollStorage.Get()
	poll.Lock.RLock()
	err := json.NewEncoder(w).Encode(poll)
	poll.Lock.RUnlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("POLL [GET] :: Error:", err)
	}
}

// PUT
func (ctrl *PollController) handleCreateOrUpdatePoll(w http.ResponseWriter, r *http.Request) {
	poll := &model.Poll{}
	err := json.NewDecoder(r.Body).Decode(poll)
	if err != nil {
		// Неверный формат опроса
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("POLL [PUT] :: Error:", err)
		return
	}

	upd := App.PollStorage.CreateOrUpdate(poll)
	if upd != nil {
		// Опрос был обновлен
		func() {
			upd.Lock.RLock()
			defer upd.Lock.RUnlock()
			App.PubSub.NotifyAll(About.UpdatedPoll(poll))
		}()
		log.Println("POLL [PUT] :: Poll update")
	}
}

// DELETE
func (ctrl *PollController) handleDeletePoll(w http.ResponseWriter, r *http.Request) {
	del := App.PollStorage.Delete()
	if del != nil {
		// Опрос не был удален ранее
		App.PubSub.NotifyAll(About.UpdatedPoll(nil))
		log.Println("POLL [DELETE] :: Poll delete")
	}
}
