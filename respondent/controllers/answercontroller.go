package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/VeselovAlex/KtoZa/model"
)

type AnswerListener interface {
	OnNewAnswerSet(model.AnswerSet)
}

// AnswerController осуществляет прием ответов
type AnswerController struct {
	listeners []AnswerListener

	lock      sync.RWMutex
	validator Validator
}

func (ctrl *AnswerController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metod not allowed. Use method POST instead", http.StatusMethodNotAllowed)
		return
	}
	var answers model.AnswerSet
	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		log.Println("ANSWER", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверка ответов
	valid := func() bool {
		ctrl.lock.RLock()
		defer ctrl.lock.RUnlock()
		if ctrl.validator.IsValid(answers) {
			return true
		}
		return false
	}()

	if valid {
		w.Write([]byte("true"))
		ctrl.notifyListeners(answers)
	} else {
		w.Write([]byte("false"))
	}
}

func (ctrl *AnswerController) notifyListeners(ans model.AnswerSet) {
	for _, listener := range ctrl.listeners {
		listener.OnNewAnswerSet(ans)
	}
}

// NewAnswerController создает новый экземпляр контроллера ответов
func NewAnswerController(listeners ...AnswerListener) *AnswerController {
	ctrl := &AnswerController{
		listeners: listeners,
	}
	return ctrl
}

func (ctrl *AnswerController) OnPollUpdate(poll *model.Poll) {
	ctrl.lock.Lock()
	defer ctrl.lock.Unlock()
	if poll != nil {
		ctrl.validator = NewValidatorFor(poll)
	}
}
