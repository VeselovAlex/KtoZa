package main

import (
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/VeselovAlex/KtoZa/model"
)

var StatisticsStorage = NewDummyStatisticsStorage()

func NewDummyStatisticsStorage() *dummyStatStorage {
	PollStorage.lock.RLock()
	defer PollStorage.lock.RUnlock()
	return &dummyStatStorage{
		stat: model.CreateStatisticsFor(PollStorage.poll),
	}
}

type dummyStatStorage struct {
	lock sync.RWMutex
	stat *model.Statistics
}

func (storage *dummyStatStorage) Get() *model.Statistics {
	storage.lock.RLock()
	retVal := storage.stat
	storage.lock.RUnlock()
	return retVal
}

func (storage *dummyStatStorage) WriteJSON(w io.Writer) error {
	storage.lock.RLock()
	defer storage.lock.RUnlock()
	return json.NewEncoder(w).Encode(storage.stat)
}

func (storage *dummyStatStorage) ApplyAnswerSet(set model.AnswerSet) bool {
	storage.lock.Lock()
	defer storage.lock.Unlock()
	ok := storage.stat.ApplyAnswerSet(set)
	if ok {
		storage.onUpdate()
	}
	return ok
}

func (storage *dummyStatStorage) onUpdate() {
	log.Println("Updated")
}
