package controllers

import (
	"encoding/json"
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

	stats  chan *model.Statistics
	polls  chan *model.Poll
	caches chan *model.Statistics

	conn *websocket.Conn

	lock sync.RWMutex
	stat *model.Statistics
}

// ConnectToMaster инициализирует прокси M-Сервера с данным URL
func ConnectToMaster(urlString string) {
	MasterServer = &master{
		hostUrl: urlString,
		stats:   make(chan *model.Statistics, 8),
		caches:  make(chan *model.Statistics, 8),
		polls:   make(chan *model.Poll, 8),
	}

	url := strings.Replace(MasterServer.hostUrl, "http:", "ws:", 1) + "/api/ws"
	conn, err := websocket.Dial(url, "", MasterServer.hostUrl)
	if err != nil {
		log.Fatalln("MASTER SERVER :: Unable to connect:", err)
	}
	MasterServer.conn = conn
	go MasterServer.read()
	go MasterServer.write()
}

// GetPoll запрашивает данные опроса с M-Сервера. Если данные успешно получены, возвращает
// ссылку на экземпляр опроса, в противном случае возвращает nil и ошибку
func (m *master) GetPoll() (*model.Poll, error) {
	resp, err := http.Get(m.hostUrl + "/api/poll")
	if err != nil {
		return nil, err
	}

	poll := &model.Poll{}
	err = json.NewDecoder(resp.Body).Decode(poll)
	if err != nil {
		return nil, err
	}
	return poll, nil
}

// GetPoll запрашивает данные статистики с M-Сервера. Если данные успешно получены, возвращает
// ссылку на экземпляр статистики, в противном случае возвращает nil и ошибку
func (m *master) GetStatistics() (*model.Statistics, error) {
	resp, err := http.Get(m.hostUrl + "/api/stats")
	if err != nil {
		return nil, err
	}

	stat := &model.Statistics{}
	err = json.NewDecoder(resp.Body).Decode(stat)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

// AwaitPollUpdate ожидает изменение опроса и возвращает экземпляр изменненного опроса
func (m *master) AwaitPollUpdate() *model.Poll {
	return <-m.polls
}

// AwaitStatisticsUpdate ожидает изменение статистики и возвращает экземпляр изменненной статистики
func (m *master) AwaitStatisticsUpdate() *model.Statistics {
	return <-m.stats
}

// SendAnswerCache асинхронно отправляет заданный кэш статистики
func (m *master) SendAnswerCache(cache *model.Statistics) {
	if cache == nil {
		return
	}
	cpy := &model.Statistics{}
	cache.CopyTo(cpy)
	m.caches <- cpy
}

func (m *master) write() {
	const maxBadsInRow = 3
	badsInRow := 0
	for cache := range m.caches {
		msg := common.Wrap.NewAnswerCache(cache)
		err := websocket.JSON.Send(m.conn, msg)
		if err != nil {
			badsInRow++
			log.Println("MASTER SERVER :: Connection lost:", err)
			if badsInRow >= maxBadsInRow {
				log.Fatalf("MASTER SERVER :: Server dropped after %d bad connections in row reached\n", badsInRow)
			}
			continue
		}
		badsInRow = 0
		log.Println("MASTER SERVER :: Answer cache submitted successfully")
	}
}

func (m *master) read() {
	const maxBadsInRow = 3
	badsInRow := 0
	for {
		msg := &common.EventRawMessage{}
		err := websocket.JSON.Receive(m.conn, msg)
		if err != nil {
			badsInRow++
			log.Println("MASTER SERVER :: Bad connection:", err)
			if badsInRow >= maxBadsInRow {
				log.Fatalf("MASTER SERVER :: Server dropped after %d bad connections in row reached\n", badsInRow)
			}
			time.Sleep(time.Second)
			continue
		}
		badsInRow = 0
		switch msg.Event {
		case common.EventUpdatedPoll:
			poll := &model.Poll{}
			err = json.Unmarshal(msg.Data, poll)
			if err != nil {
				log.Println("MASTER SERVER :: Bad poll message, skip")
			} else {
				m.polls <- poll
				log.Println("MASTER SERVER :: Got updated poll")
			}
		case common.EventUpdatedStatistics:
			stat := &model.Statistics{}
			err = json.Unmarshal(msg.Data, stat)
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
