package controllers

import (
	"encoding/json"
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

func (ctrl *PollController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metod not allowed. Use method GET instead", http.StatusMethodNotAllowed)
		return
	}

	encoder := json.NewEncoder(w)
	// Оборачиваем в замыкание для гарантии снятия блокировки
	err := func() error {
		ctrl.lock.RLock()
		defer ctrl.lock.RUnlock()
		return encoder.Encode(ctrl.poll)
	}()
	if err != nil {
		log.Println("POLL CONTROLLER :: Error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Println("POLL CONTROLLER :: Success")
	}
}

func (ctrl *PollController) notifyListeners(poll *model.Poll) {
	for _, listener := range ctrl.listeners {
		listener.OnPollUpdate(poll)
	}
}

func NewTestPollCtrl(listeners ...PollListener) *PollController {
	ctrl := &PollController{
		poll: &model.Poll{
			Title: "Test poll",
		},
		listeners: listeners,
	}
	ctrl.notifyListeners(ctrl.poll)
	return ctrl
}
