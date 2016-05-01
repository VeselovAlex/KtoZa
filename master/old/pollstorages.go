package main

import (
	"sync"
	"time"

	common "github.com/VeselovAlex/KtoZa"
	"github.com/VeselovAlex/KtoZa/model"
)

// NewMasterPollStorage возвращает интерфейс хранилища опросов для отладки
func NewMasterPollStorage() common.PollStorage {
	return &masterPollStorage{
		poll: &model.Poll{
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

type masterPollStorage struct {
	// Грубая блокировка
	lock sync.RWMutex
	poll *model.Poll
}

func (st *masterPollStorage) Get() *model.Poll {
	st.lock.RLock()
	poll := st.poll
	st.lock.RUnlock()
	return poll
}

func (st *masterPollStorage) CreateOrUpdate(poll *model.Poll) *model.Poll { /*
		os.Mkdir("data", 0700)
		file, err := os.Create("data/poll.json")
		if err != nil {
			log.Println("POLLSTORAGE :: Error:", err)
			return nil
		}
		if file != nil {
			defer file.Close()
		}*/
	st.lock.Lock()
	defer st.lock.Unlock()
	/*
		err = json.NewEncoder(file).Encode(poll)
		if err != nil {
			log.Println("POLLSTORAGE :: Error:", err)
			return nil
		}*/
	st.poll = poll
	return poll
}

func (st *masterPollStorage) Delete() *model.Poll {
	st.lock.Lock()
	poll := st.poll
	st.lock.Unlock()
	return poll
}
