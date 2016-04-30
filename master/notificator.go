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

func (*eventMessageWrapper) UpdatedPoll(poll *model.Poll) interface{} {
	rawPoll, err := json.Marshal(poll)
	if err != nil {
		log.Println("Message wrapper :: bad JSON marshalling,", err)
		return nil
	}
	return eventMessage{
		Event: "poll-update",
		Data:  json.RawMessage(rawPoll),
	}
}
