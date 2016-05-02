package controllers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/VeselovAlex/KtoZa/model"

	"github.com/satori/go.uuid"
)

type sessionPool struct {
	lock sync.RWMutex
	data map[string]bool
}

// New создает и регистрирует в хранилище новый ключ
func (s *sessionPool) New() string {
	key := uuid.NewV4().String()
	s.lock.Lock()
	s.data[key] = true
	s.lock.Unlock()
	return key
}

// Contains проверяет наличие сессионного ключа в хранилище
func (s *sessionPool) Contains(key string) bool {
	s.lock.RLock()
	res, ok := s.data[key]
	s.lock.RUnlock()
	return ok && res
}

// Delete удаляет сессионный ключ из хранилища
// Рекомендуется проверять значение ключа перед удалением
func (s *sessionPool) Delete(key string) {
	s.lock.RLock()
	s.data[key] = false
	s.lock.RUnlock()
}

const regKeyCookieName = "reg-key"

type SessionController struct {
	sessionPool *sessionPool

	lock    sync.RWMutex
	expires time.Time
}

func NewSessionController() *SessionController {
	return &SessionController{
		sessionPool: &sessionPool{
			data: make(map[string]bool, 128),
		},
	}
}

func (ctrl *SessionController) OnPollUpdate(poll *model.Poll) {
	ctrl.lock.Lock()
	if poll != nil {
		ctrl.expires = poll.Events.EndAt
	}
	ctrl.lock.Unlock()
}

func (ctrl *SessionController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctrl.handleCheckRegistration(w, r)
	case http.MethodPost:
		ctrl.handleRegister(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET
func (ctrl *SessionController) handleRegister(w http.ResponseWriter, r *http.Request) {
	key := ctrl.sessionPool.New()
	ctrl.lock.RLock()
	expires := ctrl.expires
	ctrl.lock.RUnlock()
	cookie := &http.Cookie{
		Name:    regKeyCookieName,
		Value:   key,
		Expires: expires,
	}
	http.SetCookie(w, cookie)
	w.Write([]byte("registered"))
}

// POST
func (ctrl *SessionController) handleCheckRegistration(w http.ResponseWriter, r *http.Request) {
	regCookie, err := r.Cookie(regKeyCookieName)
	jsonData := "false"
	if err == http.ErrNoCookie {
		w.Write([]byte(jsonData))
		return
	} else if err != nil {
		log.Println("Server error while registration check:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ctrl.sessionPool.Contains(regCookie.Value) {
		jsonData = "true"
	}
	w.Write([]byte(jsonData))
}
