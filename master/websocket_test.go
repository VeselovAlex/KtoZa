package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

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

func (n *testNotificator) Await() interface{} {
	return nil
}

func TestWebSocketConn(t *testing.T) {
	srv := httptest.NewServer(NewWebSocketPubSubController())
	origin := srv.URL
	url := strings.Replace(origin, "http:", "ws:", 1) + "/ws"

	var wg sync.WaitGroup

	var numConns = 10
	for i := 1; i <= numConns; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			conn, err := websocket.Dial(url, "", origin)
			if err != nil {
				t.Fatal("Bad websocket connection to", url)
				return
			}

			err = websocket.JSON.Send(conn, About.UpdatedPoll(nil))
			if err != nil {
				t.Error("Unable to send data")
			}

			rec := ""
			err = websocket.JSON.Receive(conn, &rec)
			if err != nil {
				t.Errorf("Client #%d::Error receiving message: %v\n", num, err)
			} else if rec != "test message" {
				t.Errorf("Client #%d::Wrong message received: expected %s, got %s\n", num, "test message", rec)
			}
			conn.Close()
		}(i)
	}

	go func() {
		for i := 0; i < numConns; i++ {
			App.PubSub.Await()
		}
	}()

	time.Sleep(100 * time.Millisecond)
	App.PubSub.NotifyAll("test message")
	wg.Wait()
	// Ожидание вывода
	time.Sleep(100 * time.Millisecond)
}
