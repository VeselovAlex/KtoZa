package main

import (
	"log"
	"net/http"
	"sync"
	"time"

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

// Delete удаляет сессионный ключ из памяти
// Рекомендуется проверять значение ключа перед удалением
func (s *sessionPool) Delete(key string) {
	s.lock.RLock()
	s.data[key] = false
	s.lock.RUnlock()
}

var SessionPool = &sessionPool{
	data: make(map[string]bool, 128),
}

const regKeyCookieName = "reg-key"

type SessionController struct{}

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

func (ctrl *SessionController) handleRegister(w http.ResponseWriter, r *http.Request) {
	key := SessionPool.New()
	expires := func() time.Time {
		PollStorage.lock.RLock()
		defer PollStorage.lock.RUnlock()
		return PollStorage.poll.Events.StartAt
	}()
	cookie := &http.Cookie{
		Name:    regKeyCookieName,
		Value:   key,
		Expires: expires,
	}
	http.SetCookie(w, cookie)
	w.Write([]byte("registered"))
}

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

	if SessionPool.Contains(regCookie.Value) {
		jsonData = "true"
	}
	w.Write([]byte(jsonData))
}
