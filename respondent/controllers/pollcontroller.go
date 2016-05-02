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

// NewPollController создает экземпляр контроллера опроса
func NewPollController(listeners ...PollListener) *PollController {
	poll, err := MasterServer.GetPoll()
	if err != nil {
		log.Fatalln("POLL CONTROLLER :: Unable to get data from server:", err)
	}
	log.Println("POLL CONTROLLER :: Got poll from master")
	ctrl := &PollController{
		listeners: listeners,
		poll:      poll,
	}

	// Для отладки, удалить после
	//now := time.Now()
	//ctrl.poll.Events.RegistrationAt = now.Add(5 * time.Second)
	//ctrl.poll.Events.StartAt = now.Add(10 * time.Second)
	//ctrl.poll.Events.EndAt = now.Add(30 * time.Second)
	//
	ctrl.notifyListeners(ctrl.poll)
	go ctrl.listenToMaster()
	return ctrl
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

func (ctrl *PollController) listenToMaster() {
	updWith := func(poll *model.Poll) {
		ctrl.lock.Lock()
		defer ctrl.lock.Unlock()
		ctrl.poll = poll
		ctrl.notifyListeners(ctrl.poll)
	}
	for {
		poll := MasterServer.AwaitPollUpdate()
		log.Println("POLL CONTROLLER :: Got poll update")
		updWith(poll)
	}
}
