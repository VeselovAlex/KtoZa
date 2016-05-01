package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/VeselovAlex/KtoZa/model"
)

// PollListener представляет интерфейс объекта,
// отслеживающего изменения опроса
type PollListener interface {
	OnPollUpdate(*model.Poll)
}

// PollController обрабатывает запросы, связанные с данными опроса
type PollController struct {
	listeners []PollListener

	lock sync.RWMutex
	poll *model.Poll
}

func NewPollController(listeners ...PollListener) *PollController {
	poll, err := Storage.ReadPoll()
	if err != nil {
		if err == os.ErrNotExist {
			poll = &model.Poll{}
			err = Storage.WritePoll(poll)
			if err != nil {
				log.Fatalln("POLL CONTROLLER :: Unable to read poll data:", err)
			}
		} else {
			log.Fatalln("POLL CONTROLLER :: Unable to read poll data:", err)
		}
	}

	ctrl := &PollController{
		poll:      poll,
		listeners: listeners,
	}
	ctrl.notifyListeners(ctrl.poll)
	return ctrl
}

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
		log.Printf("POLL CONTROLLER :: [%s] %s\n", r.Method, errMsg)
	}
}

// GET
func (ctrl *PollController) handleGetPoll(w http.ResponseWriter, r *http.Request) {
	ctrl.lock.RLock()
	defer ctrl.lock.RUnlock()
	err := json.NewEncoder(w).Encode(ctrl.poll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("POLL CONTROLLER :: [GET] Error:", err)
	}
}

// PUT
func (ctrl *PollController) handleCreateOrUpdatePoll(w http.ResponseWriter, r *http.Request) {
	poll := &model.Poll{}
	err := json.NewDecoder(r.Body).Decode(poll)
	if err != nil {
		// Неверный формат опроса
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("POLL CONTROLLER :: [PUT] Error:", err)
		return
	}

	ctrl.lock.Lock()
	defer ctrl.lock.Unlock()
	ctrl.poll = poll
	ctrl.notifyListeners(poll)
	log.Println("POLL CONTROLLER :: [PUT] Poll update")
}

// DELETE
func (ctrl *PollController) handleDeletePoll(w http.ResponseWriter, r *http.Request) {
	ctrl.lock.Lock()
	defer ctrl.lock.Unlock()
	if ctrl.poll != nil {
		// Опрос не был удален ранее
		ctrl.poll = nil
		ctrl.notifyListeners(nil)
		log.Println("POLL CONTROLLER :: [DELETE] Poll delete")
	}
}

func (ctrl *PollController) notifyListeners(poll *model.Poll) {
	for _, listener := range ctrl.listeners {
		listener.OnPollUpdate(poll)
	}
}
