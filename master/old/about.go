package main

import (
	"encoding/json"
	"log"

	"github.com/VeselovAlex/KtoZa/model"
)

var About *eventMessageWrapper

type eventMessageWrapper struct{}

type eventMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

func (*eventMessageWrapper) updated(evtMsg string, data interface{}) interface{} {
	raw, err := json.Marshal(data)
	if err != nil {
		log.Println("Message wrapper :: bad JSON marshalling,", err)
		return nil
	}
	return eventMessage{
		Event: evtMsg,
		Data:  json.RawMessage(raw),
	}
}

const (
	EventUpdatedPoll       = "poll-update"
	EventUpdatedStatistics = "stats-update"
	EventNewAnswerCache    = "answer-cache"
)

func (w *eventMessageWrapper) UpdatedPoll(poll *model.Poll) interface{} {
	return w.updated(EventUpdatedPoll, poll)
}

func (w *eventMessageWrapper) UpdatedStatistics(stat *model.Statistics) interface{} {
	return w.updated(EventUpdatedStatistics, stat)
}
