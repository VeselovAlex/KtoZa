package model

import (
	"encoding/json"
	"strings"
	"testing"
)

const testJSON = `{
    "title": "Простой опрос",
    "caption": "Простой опрос для отладки",
    "events": {
        "registration": "2016-04-29T00:39:30.3085295+03:00",
        "start": "2016-04-29T00:40:00.3085295+03:00",
        "end": "2016-04-29T22:39:50.3085295+03:00"
    },
    "questions": [
        {
            "text": "Вопрос 1",
            "type": "single-option",
            "options": [
                "Вариант 1",
                "Вариант 2",
                "Вариант 3"
            ]
        },
        {
            "text": "Вопрос 2",
            "type": "multi-option",
            "options": [
                "Вариант 1",
                "Вариант 2",
                "Вариант 3"
            ]
        }
    ]
}
`

func checkEquals(t *testing.T, got, expect interface{}) {
	if got != expect {
		t.Error("Bad JSON parcing: expected", expect, "got", got)
	}
}

func TestParsingJSON(t *testing.T) {
	poll := &Poll{}

	reader := strings.NewReader(testJSON)
	err := json.NewDecoder(reader).Decode(poll)

	if err != nil {
		t.Fatal("Bad JSON parsing:", err)
	}

	checkEquals(t, poll.Title, "Простой опрос")
}
