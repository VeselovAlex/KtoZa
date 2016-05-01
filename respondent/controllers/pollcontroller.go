package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

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

func (ctrl *PollController) doReadPollFrom(r io.Reader) bool {
	ctrl.lock.Lock()
	defer ctrl.lock.Unlock()
	poll := &model.Poll{}
	err := json.NewDecoder(r).Decode(poll)
	if err != nil {
		log.Println("POLL CONTROLLER :: Unable to read poll:", err)
		return false
	}
	ctrl.poll = poll
	return true
}

func NewTestPollCtrl(listeners ...PollListener) *PollController {
	src, err := os.Open(path.Join("testdata", "poll.json"))
	if err != nil {
		log.Fatalln("POLL CONTROLLER :: Initialization failed:", err)
	}
	defer src.Close()

	ctrl := &PollController{
		listeners: listeners,
	}
	ok := ctrl.doReadPollFrom(src)
	if !ok {
		log.Fatalln("POLL CONTROLLER :: Initialization failed:", err)
	}
	now := time.Now()
	ctrl.poll.Events.RegistrationAt = now.Add(5 * time.Second)
	ctrl.poll.Events.StartAt = now.Add(10 * time.Second)
	ctrl.poll.Events.EndAt = now.Add(30 * time.Second)
	ctrl.notifyListeners(ctrl.poll)
	return ctrl
}
