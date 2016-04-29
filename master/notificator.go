package main

import (
	"log"

	"github.com/VeselovAlex/KtoZa/model"
)

type Notificator interface {
	NotifyAll(interface{})
}

var Respondents *logNotificator

type logNotificator struct{}

func (*logNotificator) NotifyAll(msg interface{}) {
	about, ok := msg.(eventMessage)
	if ok {
		log.Println("NOTIFICATOR :: Received notification about", about.Event)
	} else {
		log.Println("NOTIFICATOR :: Received bad message")
	}
}

/*
type PubSub interface {
	Publish(interface{})
	Subscribe(io.ReadWriteCloser)
}

type pubSubPool struct {
	subscribe chan io.ReadWriteCloser
	publish   chan interface{}
	clients   map[io.ReadWriteCloser]bool
}

func (pubSub *pubSubPool) Subscribe(client io.ReadWriteCloser) {
	pubSub.subscribe <- client
	go func() {
		for data := range pubSub.publish {

		}
	}()
}

func (pubSub *pubSubPool) Publish(data interface{}) {

}
*/

var About *eventMessageWrapper

type eventMessageWrapper struct{}

type eventMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func (*eventMessageWrapper) UpdatedPoll(poll *model.Poll) interface{} {
	return eventMessage{
		Event: "poll-update",
		Data:  poll,
	}
}
