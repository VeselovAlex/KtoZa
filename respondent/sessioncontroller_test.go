package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	regSrv := httptest.NewServer(&SessionController{})
	body := strings.NewReader("")

	// Регистрация
	resp, err := http.Post(regSrv.URL, "text/plain", body)
	if err != nil {
		t.Fatal("Unable to connect to server:", err)
	}
	if resp.StatusCode >= 300 {
		t.Fatal("Server reponded with invalid status", resp.StatusCode)
	}

	// Проверяем установку cookie
	regKey := "none"
	cookies := resp.Cookies()
	var regCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == regKeyCookieName {
			regKey = cookie.Value
			regCookie = cookie
		}
	}
	t.Log("Registration key:", regKey)
	if regKey == "none" {
		t.Error("No registration key cookie")
	}

	// Проверка регистрации
	body = strings.NewReader("")
	req, err := http.NewRequest(http.MethodGet, regSrv.URL, body)
	req.AddCookie(regCookie)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("Unable to connect to server:", err)
	}
	if resp.StatusCode >= 300 {
		t.Fatal("Server reponded with invalid status", resp.StatusCode)
	}
	var isRegistered bool
	err = json.NewDecoder(resp.Body).Decode(&isRegistered)
	if err != nil {
		t.Error("Bad response data parse:", err)
	}
	if !isRegistered {
		t.Error("Bad registration. Not registered")
	}
}
