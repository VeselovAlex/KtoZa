package main

import (
	"bytes"
	"io"
	"log"

	"github.com/VeselovAlex/KtoZa/model"
)

type PubSub interface {
	NotifyAll(interface{})
	Subscribe(io.ReadWriteCloser)
	Await() interface{}
}

var Respondents PubSub = NewPubSubPool()

type logNotificator struct{}

func (*logNotificator) NotifyAll(msg interface{}) {
	about, ok := msg.(eventMessage)
	if ok {
		log.Println("NOTIFICATOR :: Received notification about", about.Event)
	} else {
		log.Println("NOTIFICATOR :: Received bad message")
	}
}

func (*logNotificator) Subscribe(client io.ReadWriteCloser) {
	log.Println("NOTIFICATOR :: New client")
	client.Close()
}
func (*logNotificator) Await() interface{} {
	return nil
}

type pubSubPool struct {
	subscribe chan *connection
	delete    chan *connection
	publish   chan interface{}
	read      chan interface{}
	clients   map[*connection]bool
}

type connection struct {
	send chan interface{}
	conn io.ReadWriteCloser
}

func (pubSub *pubSubPool) Subscribe(client io.ReadWriteCloser) {
	c := &connection{
		send: make(chan interface{}),
		conn: client,
	}
	pubSub.subscribe <- c
	defer func() {
		pubSub.delete <- c
	}()

	// Отправка сообщения клиенту
	go func() {
		for {
			data := <-c.send
			err := App.ResponseEncoder.ToResponseWriter(c.conn, data)
			if err != nil {
				pubSub.delete <- c
				log.Println("Connection lost:", err)
				break
			}
		}
		c.conn.Close()
	}()

	// Получение сообщений от клиента
	var err error
	for {
		var buffer []byte
		_, err = c.conn.Read(buffer)
		if err == nil {
			msg := &eventMessage{}
			err = App.RequestDecoder.FromRequest(bytes.NewBuffer(buffer), msg)
			if err == nil {
				pubSub.read <- msg
			} else {
				log.Println("Bad message:", err)
			}
		} else {
			break
		}
	}
	c.conn.Close()
	log.Println("Connection lost:", err)
}

func (pubSub *pubSubPool) NotifyAll(data interface{}) {
	pubSub.publish <- data
}

func (pubSub *pubSubPool) Await() interface{} {
	return <-pubSub.read
}

func (pubSub *pubSubPool) run() {
	log.Println("PubSubPool is now running")
	for {
		select {
		case c := <-pubSub.delete:
			c.conn.Close()
			delete(pubSub.clients, c)
			close(c.send)
		case c := <-pubSub.subscribe:
			pubSub.clients[c] = true
		case data := <-pubSub.publish:
			for c := range pubSub.clients {
				select {
				case c.send <- data:
					// Отправка сообщения
				default:
					delete(pubSub.clients, c)
					close(c.send)
				}
			}
		}
	}
}

func NewPubSubPool() PubSub {
	res := &pubSubPool{
		subscribe: make(chan *connection, 16),
		delete:    make(chan *connection, 16),
		publish:   make(chan interface{}, 16),
		read:      make(chan interface{}, 16),
		clients:   make(map[*connection]bool, 16),
	}
	go res.run()
	return res
}

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
