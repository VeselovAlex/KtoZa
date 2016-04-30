package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/websocket"
)

type testNotificator struct {
	t *testing.T
}

func (n *testNotificator) NotifyAll(msg interface{}) {
	about, ok := msg.(eventMessage)
	if ok {
		n.t.Log("NOTIFICATOR :: Received notification about", about.Event)
	} else {
		n.t.Error("NOTIFICATOR :: Received bad message")
	}
}

func (n *testNotificator) Subscribe(client io.ReadWriteCloser) {
	n.t.Log("NOTIFICATOR :: New client")
	client.Close()
}

func TestWebSocketConn(t *testing.T) {
	old := Respondents
	defer func() {
		Respondents = old
	}()
	Respondents = &testNotificator{t}
	srv := httptest.NewServer(websocket.Handler(handleWebSocketConnection))
	origin := srv.URL
	url := strings.Replace(origin, "http:", "ws:", 1) + "/ws"
	_, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Error("Bad websocket connection to", url)
	}
}
