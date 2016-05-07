package controllers

import (
	"net/http"
	"time"

	"github.com/satori/go.uuid"
)

type UserControl struct {
	authCookie string
	pwd        string
}

func NewUserControl(pwd string) *UserControl {
	return &UserControl{
		authCookie: uuid.NewV4().String(),
		pwd:        pwd,
	}
}

// CheckAuth проверяет авторизацию пользователя на основе cookies
func (uc *UserControl) CheckAuth(w http.ResponseWriter, r *http.Request) bool {
	const cookieName = "MSESSION"
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value != uc.authCookie {
		//http.Error(w, "Not authorized", http.StatusForbidden)
		return false
	}
	return true
}

// DoAuthorization выполяет авторизацию пользователя на основе
// переданного в качестве параметра запроса пароля
func (uc *UserControl) DoAuthorization(w http.ResponseWriter, r *http.Request) {
	const cookieName = "MSESSION"
	gotPwd := r.FormValue("p")
	if gotPwd == uc.pwd {
		http.SetCookie(w, &http.Cookie{
			Name:    cookieName,
			Value:   uc.authCookie,
			Expires: time.Now().Add(24 * time.Hour),
		})
	}
}

// Authorized оборачивает http.Handler для проверки авторизации при запросе к серверу
// Methods -- список методов http, для которых требуется авторизация
func (uc *UserControl) Authorized(h http.Handler, methods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range methods {
			if r.Method == method {
				authorized := uc.CheckAuth(w, r)
				if !authorized {
					//return
				}
			}
		}
		h.ServeHTTP(w, r)
	})
}

func (uc *UserControl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		uc.CheckAuth(w, r)
	case http.MethodPost:
		uc.DoAuthorization(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
