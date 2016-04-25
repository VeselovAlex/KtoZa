//Package KtoZa contains core declarations for poll service
package KtoZa

import (
	"time"
)

type pollHeader struct {
	Title    string    `json:"title"`
	Caption  string    `json:"caption"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
}

type pollQuestion struct {
	Text    string       `json:"text"`
	Options []pollOption `json:"options,omitempty"`
}

type pollOption struct {
	Text string `json:"text"`
}

// Poll represents basic poll structure
type Poll struct {
	pollHeader
	Questions []pollQuestion `json:"questions"`
}

func GetDummyPoll() *Poll {
	return &Poll{
		pollHeader{
			Title:    "Dummy poll",
			StartsAt: time.Now().Add(30 * time.Second),
			EndsAt:   time.Now().Add(time.Minute),
		},
		nil,
	}
}
