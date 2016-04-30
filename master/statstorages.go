package main

import (
	common "github.com/VeselovAlex/KtoZa"
	"github.com/VeselovAlex/KtoZa/model"
)

func NewMasterStatisticsStorage() common.StatisticsStorage {
	return &masterStatStorage{
		stat: model.CreateStatisticsFor(App.PollStorage.Get()),
	}
}

type masterStatStorage struct {
	stat *model.Statistics
}

func (st *masterStatStorage) Get() *model.Statistics {
	return st.stat
}

func (st *masterStatStorage) CreateOrUpdate(stat *model.Statistics) *model.Statistics {
	st.stat = stat
	return st.stat
}

func (st *masterStatStorage) Delete() *model.Statistics {
	old := st.stat
	st.stat = nil
	return old
}
