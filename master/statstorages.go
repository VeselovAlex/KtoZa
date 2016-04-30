package main

import (
	"sync"

	common "github.com/VeselovAlex/KtoZa"
	"github.com/VeselovAlex/KtoZa/model"
)

// NewMasterStatisticsStorage возвращет экземпляр хранилища статистики для M-сервера
func NewMasterStatisticsStorage() common.StatisticsStorage {
	return &masterStatStorage{
		stat: model.CreateStatisticsFor(App.PollStorage.Get()),
	}
}

type masterStatStorage struct {
	// Грубая блокировка
	lock sync.RWMutex
	stat *model.Statistics
}

func (st *masterStatStorage) Get() *model.Statistics {
	st.lock.RLock()
	stat := st.stat
	st.lock.RUnlock()
	return st.stat
}

func (st *masterStatStorage) CreateOrUpdate(stat *model.Statistics) *model.Statistics {
	st.lock.Lock()
	defer st.lock.Unlock()
	joined := st.stat.JoinWith(stat)
	if joined {
		return stat
	}
	return nil
}

func (st *masterStatStorage) Delete() *model.Statistics {
	st.lock.Lock()
	old := st.stat
	st.stat = nil
	st.lock.Unlock()
	return old
}
