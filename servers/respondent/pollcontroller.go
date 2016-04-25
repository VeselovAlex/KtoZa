package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	core "KtoZa"

	"github.com/satori/go.uuid"
	ws "golang.org/x/net/websocket"
)

type pollController struct {
	secret string
	data   *core.Poll

	lockState sync.RWMutex
	closed    bool

	lockBroadcast sync.Mutex
	condSend      *sync.Cond
}

func (c *pollController) Run() {
	c.condSend = sync.NewCond(&c.lockBroadcast)

	delay := c.data.StartsAt.Sub(time.Now())
	<-time.After(delay)
	c.lockState.Lock()
	c.closed = true
	c.lockState.Unlock()
	c.condSend.Broadcast()
}

func (c *pollController) SendSecretAndClose(conn *ws.Conn) {
	if conn == nil {
		return
	}

	if c.IsActive() {
		_, err := fmt.Fprint(conn, c.waitForSecret())
		if err != nil {
			log.Println("Error sending:", err)
			return
		}
	}
	conn.Close()
}

// Unsafe function
// KISS
func (c *pollController) waitForSecret() string {
	c.lockBroadcast.Lock()
	c.condSend.Wait()
	secret := c.secret
	c.lockBroadcast.Unlock()
	return secret
}

func (c *pollController) IsActive() bool {
	c.lockState.RLock()
	defer c.lockState.RUnlock()
	return !c.closed
}

func (c *pollController) Get() *core.Poll {
	return c.data
}

func (c *pollController) Init() {
	PollController.data = Master.GetPoll()
	PollController.secret = uuid.NewV4().String()
}

var PollController = pollController{}
