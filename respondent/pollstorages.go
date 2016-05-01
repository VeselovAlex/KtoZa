package main

import (
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/VeselovAlex/KtoZa/model"
)

var PollStorage = NewDummyPollStorage()

// NewDummyPollStorage возвращает интерфейс хранилища опросов для отладки
/*func NewDummyPollStorage() common.PollStorage {
	resp, err := http.Get("http://localhost:8888/api/poll")
	if err != nil {
		log.Fatalln("Bad master connection")
	}
	poll := &model.Poll{}
	err = json.NewDecoder(resp.Body).Decode(poll)
	if err != nil {
		log.Fatalln("Bad master response")
	}
	log.Println("got poll", poll.Title)
	return &dummyPollStorage{poll}
}*/

func NewDummyPollStorage() *dummyPollStorage {
	return &dummyPollStorage{
		poll: &model.Poll{
			Title: "Dummy",
			Questions: []model.Question{
				model.Question{
					Type:    model.TypeSingleOptionQuestion,
					Text:    "Select 0 or 1",
					Options: []string{"0", "1"},
				},
			},
			Events: model.EventTimings{
				EndAt: time.Now(),
			},
		},
	}
}

type dummyPollStorage struct {
	lock sync.RWMutex
	poll *model.Poll
}

func (storage *dummyPollStorage) Get() *model.Poll {
	storage.lock.RLock()
	retVal := storage.poll
	storage.lock.RUnlock()
	return retVal
}

func (storage *dummyPollStorage) WriteJSON(w io.Writer) error {
	storage.lock.RLock()
	defer storage.lock.RUnlock()
	return json.NewEncoder(w).Encode(storage.poll)
}
