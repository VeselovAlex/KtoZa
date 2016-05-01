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

var poll = &Poll{}

func TestParsingJSON(t *testing.T) {
	reader := strings.NewReader(testJSON)
	err := json.NewDecoder(reader).Decode(poll)

	if err != nil {
		t.Fatal("Bad JSON parsing:", err)
	}

	checkEquals(t, poll.Title, "Простой опрос")
}

func TestAnswerVerification(t *testing.T) {
	good := []AnswerSet{
		AnswerSet([]Answer{Answer([]int{}), Answer([]int{})}),        // Пустой ответ
		AnswerSet([]Answer{Answer([]int{0}), Answer([]int{1})}),      // По одному варианту
		AnswerSet([]Answer{Answer([]int{1}), Answer([]int{0, 2})}),   // Несколько вариантов в multi-option
		AnswerSet([]Answer{Answer([]int{2}), Answer([]int{})}),       // Пустой multi-option
		AnswerSet([]Answer{Answer([]int{}), Answer([]int{2, 1, 0})}), // Обратный порядок в multi-option
	}

	for _, ans := range good {
		if !poll.IsValidAnswerSet(ans) {
			t.Error("False negative validation of answer set", ans)
		}
	}

	bad := []AnswerSet{
		AnswerSet([]Answer{Answer([]int{})}),                                // Не все вопросы
		AnswerSet([]Answer{Answer([]int{1, 2}), Answer([]int{0, 2})}),       // Несколько вариантов в single-option
		AnswerSet([]Answer{Answer([]int{1, 2}), Answer([]int{0, 2, 0, 3})}), // Слишком много вариантов в single-option
		AnswerSet([]Answer{Answer([]int{2}), Answer([]int{0, 10})}),         // Большое значение в multi-option
		AnswerSet([]Answer{Answer([]int{10}), Answer([]int{2, 1, 0})}),      // Большое значение в single-option
		AnswerSet([]Answer{Answer([]int{2}), Answer([]int{0, -1})}),         // Отрицательное значение в multi-option
		AnswerSet([]Answer{Answer([]int{-1}), Answer([]int{2, 1, 0})}),      // Отрицательное значение в single-option
	}

	for _, ans := range bad {
		if poll.IsValidAnswerSet(ans) {
			t.Error("False positive validation of answer set", ans)
		}
	}
}
