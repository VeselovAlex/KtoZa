package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VeselovAlex/KtoZa/model"
)

const testRespondString = "Test"

func TestGetPoll(t *testing.T) {
	srv := httptest.NewServer(&PollController{})
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal("Unable to start test server", err)
	}

	poll := &model.Poll{}
	err = json.NewDecoder(res.Body).Decode(poll)
	if err != nil {
		t.Fatal("Unable to decode poll from response:", err)
	}
}
