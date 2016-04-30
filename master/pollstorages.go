package main

import (
	"log"
	"time"

	common "github.com/VeselovAlex/KtoZa"
	"github.com/VeselovAlex/KtoZa/model"
)

// NewDummyPollStorage возвращает интерфейс хранилища опросов для отладки
func NewDummyPollStorage() common.PollStorage {
	return &dummyPollStorage{
		&model.Poll{
			Title:   "Dummy poll (master)",
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
	App.PubSub.NotifyAll(About.UpdatedPoll(poll))
	return poll
}

func (st *dummyPollStorage) Delete() *model.Poll {
	log.Println("Dummy poll storage Delete()")
	return st.Get()
}
