package main

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type wsPubSubPool struct {
	subscribe chan *connection
	delete    chan *connection
	publish   chan interface{}
	read      chan interface{}
	clients   map[*connection]bool
}

func newWebSocketPubSub() *wsPubSubPool {
	res := &wsPubSubPool{
		subscribe: make(chan *connection, 16),
		delete:    make(chan *connection, 16),
		publish:   make(chan interface{}, 16),
		read:      make(chan interface{}, 16),
		clients:   make(map[*connection]bool, 16),
	}
	go res.run()
	return res
}

func (pubSub *wsPubSubPool) Subscribe(client io.ReadWriteCloser) {
	socket, ok := client.(*websocket.Conn)
	if ok {
		pubSub.socketSubscribe(socket)
	} else {
		log.Println("Not websocket connection, aborted")
	}
}

func (pubSub *wsPubSubPool) socketSubscribe(socket *websocket.Conn) {
	c := &connection{
		send:   make(chan interface{}),
		socket: socket,
		pool:   pubSub,
	}
	pubSub.subscribe <- c
	defer func() {
		pubSub.delete <- c
	}()

	// Отправка сообщения клиенту
	go c.write()

	// Получение сообщений от клиента
	c.read()
}

func (pubSub *wsPubSubPool) NotifyAll(data interface{}) {
	pubSub.publish <- data
}

func (pubSub *wsPubSubPool) Await() interface{} {
	return <-pubSub.read
}

func (pubSub *wsPubSubPool) run() {
	log.Println("wsPubSubPool is now running")
	for {
		select {
		case c := <-pubSub.delete:
			// Удаление клиента
			delete(pubSub.clients, c)
			close(c.send)
		case c := <-pubSub.subscribe:
			// Добавление клиента
			pubSub.clients[c] = true
		case data := <-pubSub.publish:
			for c := range pubSub.clients {
				select {
				case c.send <- data:
					// Отправка сообщения
				default:
					// Не могу отправить сообщение
					delete(pubSub.clients, c)
				}
			}
		}
	}
}

type connection struct {
	send   chan interface{}
	socket *websocket.Conn
	pool   *wsPubSubPool
}

func (c *connection) read() {
	for {
		msg := &eventMessage{}
		err := websocket.JSON.Receive(c.socket, msg)
		if err != nil {
			log.Println("WEBSOCKET PUB/SUB :: Bad connection response,", err, "[closing]")
			break
		}
		log.Println("Got", msg)
	}
	c.socket.Close()
	log.Println("Closed (reading)")
}

func (c *connection) write() {
	for msg := range c.send {
		err := websocket.JSON.Send(c.socket, msg)
		if err != nil {
			break
		}
		log.Println("Sent ::", msg)
	}
	c.socket.Close()
	log.Println("Closed (writing)")
}

func NewWebSocketPubSubController() http.Handler {
	if App.PubSub == nil {
		App.PubSub = newWebSocketPubSub()
	}
	return websocket.Handler(handleWebSocketConnection)
}

func handleWebSocketConnection(conn *websocket.Conn) {
	App.PubSub.Subscribe(conn)
}
