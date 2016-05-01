package controllers

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"sync"

	"github.com/VeselovAlex/KtoZa/model"
)

var Storage *storage

type storage struct {
	dataPath string

	pollLock sync.Mutex
	statLock sync.Mutex
}

func LoadFileSystemStorage(path string) {
	Storage = &storage{dataPath: path}
}

func (st *storage) ReadPoll() (*model.Poll, error) {
	st.pollLock.Lock()
	defer st.pollLock.Unlock()
	path := path.Join(st.dataPath, "poll.json")
	src, err := os.Open(path)
	if src != nil {
		defer src.Close()
	}
	if err != nil {
		return nil, err
	}

	poll := &model.Poll{}
	err = json.NewDecoder(src).Decode(poll)
	if err != nil {
		return nil, err
	}
	return poll, nil
}

func (st *storage) WritePoll(poll *model.Poll) error {
	st.pollLock.Lock()
	defer st.pollLock.Unlock()
	path := path.Join(st.dataPath, "poll.json")
	src, err := os.Create(path)
	if src != nil {
		defer src.Close()
	}
	if err != nil {
		return err
	}

	err = json.NewEncoder(src).Encode(poll)
	if err != nil {
		return err
	}
	return nil
}

func (st *storage) ReadStatistics() (*model.Statistics, error) {
	st.statLock.Lock()
	defer st.statLock.Unlock()
	path := path.Join(st.dataPath, "stat.json")
	src, err := os.Open(path)
	if src != nil {
		defer src.Close()
	}
	if err != nil {
		return nil, err
	}

	stat := &model.Statistics{}
	err = json.NewDecoder(src).Decode(stat)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (st *storage) WriteStatistics(stat *model.Statistics) error {
	st.statLock.Lock()
	defer st.statLock.Unlock()
	path := path.Join(st.dataPath, "stat.json")
	src, err := os.Create(path)
	if src != nil {
		defer src.Close()
	}
	if err != nil {
		return err
	}

	err = json.NewEncoder(src).Encode(stat)
	if err != nil {
		return err
	}
	return nil
}

type StorageUpdateListener struct{}

func NewStorageUpdateListener() *StorageUpdateListener {
	return &StorageUpdateListener{}
}

func (l *StorageUpdateListener) OnPollUpdate(poll *model.Poll) {
	err := Storage.WritePoll(poll)
	if err != nil {
		log.Println("STORAGE :: Unable to persist poll:", err)
	}
}

func (l *StorageUpdateListener) OnStatisticsUpdate(stat *model.Statistics) {
	err := Storage.WriteStatistics(stat)
	if err != nil {
		log.Println("STORAGE :: Unable to persist statistics:", err)
	}
}
