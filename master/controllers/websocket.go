// Александр Веселов <veselov143@gmail.com>
// СПбГУ, Математико-механический факультет, гр. 442
// Май, 2016 г.

// websocket.go содержит реализацию контроллера WebSocket-соединений
// M-сервера системы KtoZa

package controllers

import (
	"io"
	"log"
	"net/http"

	"github.com/VeselovAlex/KtoZa/model"

	common "github.com/VeselovAlex/KtoZa"
	"golang.org/x/net/websocket"
)

// Respondents представляет глобальный пул WebSocket-соединений
var Respondents = newWebSocketPubSub()

// WebSocketController осуществляет обработку WebSocket-соединений
type WebSocketController struct {
	http.Handler
}

// NewWebSocketController создает экземпляр WebSocket-контроллера
func NewWebSocketController() *WebSocketController {
	return &WebSocketController{
		websocket.Handler(handleWebSocketConnection),
	}
}

// handleWebSocketConnection добавляет соединение к
// пулу соединений при подключении
func handleWebSocketConnection(conn *websocket.Conn) {
	Respondents.Subscribe(conn)
}

// OnPollUpdate осуществляет оповещение подключенных
// клиентов об изменении статистики
func (ctrl *WebSocketController) OnPollUpdate(poll *model.Poll) {
	Respondents.NotifyAll(common.Wrap.UpdatedPoll(poll))
}

// OnStatisticsUpdate осуществляет оповещение подключенных
// клиентов об изменении статистики
func (ctrl *WebSocketController) OnStatisticsUpdate(stat *model.Statistics) {
	Respondents.NotifyAll(common.Wrap.UpdatedStatistics(stat))
}

// wsPubSubPool представляет пул WebSocket-соединений
// Основано на: [Mat Ryer. Go Programming Blueprints. Ch.1]
type wsPubSubPool struct {
	// Канал новых соединений
	subscribe chan *connection
	// Канал соединений требующих удаления
	delete chan *connection
	// Канал сообщений для отправки клиентам
	publish chan interface{}
	// Канал сообщений для приема сообщений
	read    chan interface{}
	clients map[*connection]bool
}

// newWebSocketPubSub создает и инициализирует пул соединений
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

// Subscribe проверяет, что client является WebSocket-соединением
// и добавляет клиента к пулу соединений
func (pubSub *wsPubSubPool) Subscribe(client io.ReadWriteCloser) {
	socket, ok := client.(*websocket.Conn)
	if ok {
		pubSub.socketSubscribe(socket)
	} else {
		log.Println("RESPONDENTS :: Not websocket connection, aborted")
	}
}

// socketSubscribe добавляет WebSocket-соединение к пулу соединений
func (pubSub *wsPubSubPool) socketSubscribe(socket *websocket.Conn) {
	c := &connection{
		send:   make(chan interface{}),
		socket: socket,
		pool:   pubSub,
	}
	pubSub.subscribe <- c
	defer func() { pubSub.delete <- c }()

	// Отправка сообщения клиенту
	go c.write()

	// Получение сообщений от клиента
	c.read()
}

// NotifyAll рассылает сообщение всем подключенным клиентам
func (pubSub *wsPubSubPool) NotifyAll(data interface{}) {
	pubSub.publish <- data
}

// Await возращает сообщение при его получении
func (pubSub *wsPubSubPool) Await() interface{} {
	data := <-pubSub.read
	log.Println("RESPONDENTS :: Got new message")
	return data
}

// run запускает цикл обработки сообщений
func (pubSub *wsPubSubPool) run() {
	for {
		select {
		case c := <-pubSub.delete:
			delete(pubSub.clients, c)
			close(c.send)
			log.Println("RESPONDENTS :: Removed connection")
		case c := <-pubSub.subscribe:
			pubSub.clients[c] = true
			log.Println("RESPONDENTS :: New R-server connection")
		case data := <-pubSub.publish:
			for c := range pubSub.clients {
				select {
				case c.send <- data:
					// Отправка сообщения
				default:
					// Ничего не делаем, удалением занимается другой поток
				}
			}
		}
	}
}

// conneсtion представляет клиентское соединение
type connection struct {
	send   chan interface{}
	socket *websocket.Conn
	pool   *wsPubSubPool
}

// read принимает и декодирует сообщения от клиента и
// направляет их в канал read пула
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

// write отправляет сообщения из канала send удаленнному клиенту
func (c *connection) write() {
	for msg := range c.send {
		err := websocket.JSON.Send(c.socket, msg)
		if err != nil {
			break
		}
	}
	c.socket.Close()
	log.Println("RESPONDENTS :: Connection closed (writing)")
}
