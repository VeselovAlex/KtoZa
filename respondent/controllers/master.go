package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/VeselovAlex/KtoZa/model"
)

var MasterServer *master

type master struct {
	hostUrl string

	stats  chan *model.Statistics
	caches chan *model.Statistics

	lock sync.RWMutex
	stat *model.Statistics
}

func ConnectToMaster(urlString string) error {
	MasterServer = &master{
		hostUrl: urlString,
		stats:   make(chan *model.Statistics, 8),
		caches:  make(chan *model.Statistics, 8),
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
