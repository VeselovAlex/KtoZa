package main

import (
	"log"

	"github.com/VeselovAlex/KtoZa/model"
)

type StatisticsStorage interface {
	Get() *model.Statistics
	CreateOrUpdate(*model.Statistics) *model.Statistics
	Delete() *model.Statistics
}

func NewDummyStatisticsStorage() StatisticsStorage {
	return &dummyStatStorage{
		stat: model.CreateStatisticsFor(App.PollStorage.Get()),
	}
}

type dummyStatStorage struct {
	stat *model.Statistics
}

func (st *dummyStatStorage) Get() *model.Statistics {
	log.Println("Dummy statistics storage Get()")
	return st.stat
}

func (st *dummyStatStorage) CreateOrUpdate(stat *model.Statistics) *model.Statistics {
	log.Println("Dummy statistics storage CreateOrUpdate()")
	st.stat = stat
	return st.stat
}

func (st *dummyStatStorage) Delete() *model.Statistics {
	log.Println("Dummy statistics storage Delete()")
	old := st.stat
	st.stat = nil
	return old
}
