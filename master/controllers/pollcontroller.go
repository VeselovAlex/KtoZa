package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

// NewPollController создает экземпляр контроллера опросов.
// Параметр listeners -- произвольное число слушателей обновлений опроса
func NewPollController(listeners ...PollListener) *PollController {
	// Читаем опрос из хранилища
	poll, err := Storage.ReadPoll()
	if err != nil {
		log.Println("POLL CONTROLLER :: No or bad poll data:", err)
		log.Println("POLL CONTROLLER :: Using empty poll")
	}
	listeners = append(listeners, Storage)
	ctrl := &PollController{
		poll:      poll,
		listeners: listeners,
	}
	// Оповещаем слушателей об изменении опроса
	ctrl.notifyListeners(ctrl.poll)
	return ctrl
}

// ServeHTTP реализует интерфейс http.Handler и осуществляет обработку запросов
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
