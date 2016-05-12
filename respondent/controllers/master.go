// Александр Веселов <veselov143@gmail.com>
// СПбГУ, Математико-механический факультет, гр. 442
// Май, 2016 г.
//
// master.go содержит реализацию контроллера
// взаимодействия c M-сервером
package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"

	"github.com/VeselovAlex/KtoZa/model"

	common "github.com/VeselovAlex/KtoZa"
)

// MasterServer -- прокси M-Server'а
var MasterServer *master

type master struct {
	hostUrl string
	origin  string

	stats  chan *model.Statistics
	polls  chan *model.Poll
	caches chan *model.Statistics

	errWriting chan error

	lock sync.RWMutex
	conn *websocket.Conn
}

// ConnectToMaster инициализирует прокси M-Сервера с данным URL
func ConnectToMaster(urlString, origin string) {
	MasterServer = &master{
		hostUrl:    urlString,
		origin:     origin,
		stats:      make(chan *model.Statistics, 8),
		caches:     make(chan *model.Statistics, 8),
		polls:      make(chan *model.Poll, 8),
		errWriting: make(chan error),
	}

	MasterServer.tryConnect()
	go MasterServer.read()
	go MasterServer.write()
}

func (m *master) decodePoll(from []byte) (*model.Poll, error) {
	if string(from[:4]) == "null" {
		return nil, nil
	}
	poll := &model.Poll{}
	err := json.Unmarshal(from, poll)
	return poll, err
}

func (m *master) decodeStat(from []byte) (*model.Statistics, error) {
	if string(from[:4]) == "null" {
		return nil, nil
	}
	stat := &model.Statistics{}
	err := json.Unmarshal(from, stat)
	return stat, err
}

// GetPoll запрашивает данные опроса с M-Сервера. Если данные успешно получены, возвращает
// ссылку на экземпляр опроса, в противном случае возвращает nil и ошибку
func (m *master) GetPoll() (*model.Poll, error) {
	resp, err := http.Get(m.hostUrl + "/api/poll")
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return m.decodePoll(buf)
}

// GetPoll запрашивает данные статистики с M-Сервера. Если данные успешно получены, возвращает
// ссылку на экземпляр статистики, в противном случае возвращает nil и ошибку
func (m *master) GetStatistics() (*model.Statistics, error) {
	resp, err := http.Get(m.hostUrl + "/api/stats")
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return m.decodeStat(buf)
}

// AwaitPollUpdate ожидает изменение опроса и возвращает экземпляр изменненного опроса
func (m *master) AwaitPollUpdate() *model.Poll {
	return <-m.polls
}

// AwaitStatisticsUpdate ожидает изменение статистики и возвращает экземпляр изменненной статистики
func (m *master) AwaitStatisticsUpdate() *model.Statistics {
	return <-m.stats
}

// SendAnswerCache отправляет заданный кэш статистики
func (m *master) SendAnswerCache(cache *model.Statistics) error {
	if cache == nil {
		return nil
	}
	cpy := &model.Statistics{}
	cache.CopyTo(cpy)
	m.caches <- cpy
	return <-m.errWriting
}

func (m *master) tryConnect() {
	duration := time.Second
	url := strings.Replace(m.hostUrl, "http:", "ws:", 1) + "/api/ws"
	for {
		log.Println("MASTER SERVER :: Trying to connect ...")
		conn, err := websocket.Dial(url, "", m.origin)
		if err == nil {
			m.lock.Lock()
			m.conn = conn
			m.lock.Unlock()
			log.Println("MASTER SERVER :: Connected")
			return
		}
		log.Println("MASTER SERVER :: Connection failed:", err)
		time.Sleep(duration)
	}
}

func (m *master) write() {
	for cache := range m.caches {
		if cache == nil {
			continue
		}
		msg := common.Wrap.NewAnswerCache(cache)
		err := func() error {
			m.lock.RLock()
			defer m.lock.RUnlock()
			return websocket.JSON.Send(m.conn, msg)
		}()
		if err != nil {
			log.Println("MASTER SERVER :: Connection lost:", err)
		}
		m.errWriting <- err
	}
}

func (m *master) read() {
	for {
		msg := &common.EventRawMessage{}
		err := websocket.JSON.Receive(m.conn, msg)
		if err != nil {
			log.Println("MASTER SERVER :: Connection lost:", err)
			m.tryConnect()
			continue
		}
		switch msg.Event {
		case common.EventUpdatedPoll:
			poll, err := m.decodePoll(msg.Data)
			if err != nil {
				log.Println("MASTER SERVER :: Bad poll message, skip")
			} else {
				m.polls <- poll
				log.Println("MASTER SERVER :: Got updated poll")
			}
		case common.EventUpdatedStatistics:
			stat, err := m.decodeStat(msg.Data)
			if err != nil {
				log.Println("MASTER SERVER :: Bad statistics message, skip")
			} else {
				m.stats <- stat
				log.Println("MASTER SERVER :: Got updated statistics")
			}
		default:
			log.Println("MASTER SERVER :: Unsupported message, skip")
		}
	}
}
