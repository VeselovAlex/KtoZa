package controllers

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/VeselovAlex/KtoZa/model"
)

var MasterServer *master

type master struct {
	stats  chan *model.Statistics
	caches chan *model.Statistics

	lock sync.RWMutex
	stat *model.Statistics
}

func ConnectToMaster(urlString string) error {
	MasterServer = &master{
		stats:  make(chan *model.Statistics, 8),
		caches: make(chan *model.Statistics, 8),
	}
	return nil
}

func (m *master) AwaitPollUpdate() *model.Poll {
	log.Println("Dummy poll await")
	for {

	}
	return nil
}

func (m *master) GetPoll() (*model.Poll, error) {
	src, err := os.Open(path.Join("testdata", "poll.json"))
	if err != nil {
		log.Println("MASTER SERVER :: Get poll request failed:", err)
		return nil, err
	}
	defer src.Close()

	poll := &model.Poll{}
	err = json.NewDecoder(src).Decode(poll)
	if err != nil {
		log.Println("MASTER SERVER :: Get poll request decoding failed:", err)
		return nil, err
	}
	return poll, nil
}

func (m *master) AwaitStatisticsUpdate() *model.Statistics {
	log.Println("Dummy stat await")
	return <-m.stats
}

func (m *master) SendAnswerCache(cache *model.Statistics) bool {
	var stat = new(model.Statistics)
	m.lock.Lock()
	if m.stat == nil {
		m.stat = cache
	} else {
		m.stat.JoinWith(cache)
	}
	m.stat.CopyTo(stat)
	m.lock.Unlock()
	time.Sleep(time.Second)
	m.stats <- stat
	log.Println("Dummy answer cache sent")
	return true
}
