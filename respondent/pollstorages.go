package main

import (
	"log"
	"time"

	"github.com/VeselovAlex/KtoZa/model"
)

type PollStorage interface {
	Get() *model.Poll
	CreateOrUpdate(poll *model.Poll) *model.Poll
	Delete() *model.Poll
}

// NewDummyPollStorage возвращает интерфейс хранилища опросов для отладки
func NewDummyPollStorage() PollStorage {
	return &dummyPollStorage{
		&model.Poll{
			Title:   "Dummy poll",
			Caption: "Простой опрос для тестирования",
			Events: model.EventTimings{
				RegistrationAt: time.Now().Add(5 * time.Second),
				StartAt:        time.Now().Add(35 * time.Second),
				EndAt:          time.Now().Add(95 * time.Second),
			},
			Questions: []model.Question{
				model.Question{
					Text:    "Вопрос 1",
					Type:    "single-option",
					Options: []string{"Yes", "No"},
				},
				model.Question{
					Text:    "Вопрос 2",
					Type:    "multi-option",
					Options: []string{"One", "Two", "Three"},
				},
			},
		},
	}
}

type dummyPollStorage struct {
	poll *model.Poll
}

func (st *dummyPollStorage) Get() *model.Poll {
	log.Println("Dummy poll storage Get()")
	return st.poll
}

func (st *dummyPollStorage) CreateOrUpdate(poll *model.Poll) *model.Poll {
	log.Println("Dummy poll storage CreateOrUpdate()")
	st.poll = poll
	return poll
}

func (st *dummyPollStorage) Delete() *model.Poll {
	log.Println("Dummy poll storage Delete()")
	return st.Get()
}
