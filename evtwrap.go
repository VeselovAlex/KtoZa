// Александр Веселов <veselov143@gmail.com>
// СПбГУ, Математико-механический факультет, гр. 442
// Май, 2016 г.
//
// evtwrap.go содержит реализацию обертки сообщений о происходящих событиях
package controllers

import (
	"encoding/json"

	"github.com/VeselovAlex/KtoZa/model"
)

// Wrap осущечтвляет оборачивание данных в сообщения
var Wrap *eventMessageWrapper

type eventMessageWrapper struct{}

type eventMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// EventRawMessage представляет обертку сообщения, получаемого из WebSocket-канала
type EventRawMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

const (
	EventUpdatedPoll       = "poll-update"
	EventUpdatedStatistics = "stats-update"
	EventNewAnswerCache    = "answer-cache"
)

func (w *eventMessageWrapper) UpdatedPoll(poll *model.Poll) interface{} {
	return eventMessage{EventUpdatedPoll, poll}
}

func (w *eventMessageWrapper) UpdatedStatistics(stat *model.Statistics) interface{} {
	return eventMessage{EventUpdatedStatistics, stat}
}

func (w *eventMessageWrapper) NewAnswerCache(cache *model.Statistics) interface{} {
	return eventMessage{EventNewAnswerCache, cache}
}
