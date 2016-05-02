package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"

	"github.com/VeselovAlex/KtoZa/model"
)

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

func (m *master) AwaitPollUpdate() *model.Poll {
	return <-m.polls
}

func (m *master) AwaitStatisticsUpdate() *model.Statistics {
	return <-m.stats
}

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
		msg := About.NewAnswerCache(cache).(eventMessage)
		log.Printf("DEBUG :: Sent data %s\n", msg.Data)
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
		msg := &eventMessage{}
		err := websocket.JSON.Receive(m.conn, msg)
		if err != nil {
			badsInRow++
			log.Println("MASTER SERVER :: Bad connection:", err)
			if badsInRow >= maxBadsInRow {
				log.Fatalf("MASTER SERVER :: Server dropped after %d bad connections in row reached\n", badsInRow)
			}
			time.Sleep(10 * time.Second)
			continue
		}
		badsInRow = 0
		switch msg.Event {
		case EventUpdatedPoll:
			poll := &model.Poll{}
			// Внезапно данные в Base64
			raw := []byte(msg.Data)
			// Удаляем кавычки
			raw = raw[1 : len(raw)-1]
			reader := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(raw))
			err = json.NewDecoder(reader).Decode(poll)
			if err != nil {
				log.Println("MASTER SERVER :: Bad poll message, skip")
			} else {
				m.polls <- poll
				log.Println("MASTER SERVER :: Got updated poll")
			}
		case EventUpdatedStatistics:
			stat := &model.Statistics{}
			// Внезапно данные в Base64
			raw := []byte(msg.Data)
			// Удаляем кавычки
			raw = raw[1 : len(raw)-1]
			reader := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(raw))
			err = json.NewDecoder(reader).Decode(stat)
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
