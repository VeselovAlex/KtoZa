package controllers

import (
	"io"
	"log"
	"net/http"

	"github.com/VeselovAlex/KtoZa/model"

	common "github.com/VeselovAlex/KtoZa"
	"golang.org/x/net/websocket"
)

var Respondents = newWebSocketPubSub()

type WebSocketController struct {
	http.Handler
}

func NewWebSocketController() *WebSocketController {
	return &WebSocketController{
		websocket.Handler(handleWebSocketConnection),
	}
}

func handleWebSocketConnection(conn *websocket.Conn) {
	Respondents.Subscribe(conn)
}

func (ctrl *WebSocketController) OnPollUpdate(poll *model.Poll) {
	Respondents.NotifyAll(common.Wrap.UpdatedPoll(poll))
}
func (ctrl *WebSocketController) OnStatisticsUpdate(stat *model.Statistics) {
	Respondents.NotifyAll(common.Wrap.UpdatedStatistics(stat))
}

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
		log.Println("RESPONDENTS :: Not websocket connection, aborted")
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
	data := <-pubSub.read
	log.Println("DATA :: ", data)
	log.Println("RESPONDENTS :: Got new message")
	return data
}

func (pubSub *wsPubSubPool) run() {
	for {
		select {
		case c := <-pubSub.delete:
			// Удаление клиента
			delete(pubSub.clients, c)
			close(c.send)
		case c := <-pubSub.subscribe:
			pubSub.clients[c] = true
			log.Println("RESPONDENTS :: New R-server connection")
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
		raw := common.EventRawMessage{}
		err := websocket.JSON.Receive(c.socket, &raw)
		if err != nil {
			log.Println("RESPONDENTS :: Bad connection response,", err, "[closing]")
			break
		}
		c.pool.read <- raw
	}
	c.socket.Close()
	log.Println("RESPONDENTS :: Connection closed (reading)")
}

func (c *connection) write() {
	for msg := range c.send {
		err := websocket.JSON.Send(c.socket, msg)
		if err != nil {
			break
		}
		c.pool.read <- msg
	}
	c.socket.Close()
	log.Println("RESPONDENTS :: Connection closed (writing)")
}
