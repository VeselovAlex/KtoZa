package model

import (
	"sync"
	"time"
)

// Poll представляет структуру данных опроса системы KtoZa
type Poll struct {
	Title     string       `json:"title"`
	Caption   string       `json:"caption"`
	Events    EventTimings `json:"events"`
	Questions []Question   `json:"questions"`

	// Грубая блокировка опроса
	Lock sync.RWMutex `json:"-"`
}

// EventTimings представляет структуру, хранящую время до начала событий опроса:
// начала регистрации, начала приема ответов и завершения опросов
type EventTimings struct {
	RegistrationAt time.Time `json:"registration"`
	StartAt        time.Time `json:"start"`
	EndAt          time.Time `json:"end"`
}

const (
	// TypeSingleOptionQuestion - тип вопроса с выбором одного варианта ответа
	TypeSingleOptionQuestion = "single-option"
	// TypeMultiOptionQuestion - тип вопроса с выбором одного или нескольких вариантов ответа
	TypeMultiOptionQuestion = "multi-option"
)

// Question представляет структуру вопроса
type Question struct {
	Text    string   `json:"text"`
	Type    string   `json:"type"`
	Options []string `json:"options"`
}

// IsValidAnswerSet возвращает true, если заданный набор ответов соответствует опросу
func (poll *Poll) IsValidAnswerSet(set AnswerSet) bool {
	setSlice := []Answer(set)
	// Проверка длин
	if len(setSlice) != len(poll.Questions) {
		return false
	}

	for i, ans := range setSlice {
		q := poll.Questions[i]
		switch q.Type {
		case TypeSingleOptionQuestion:
			if !poll.isValidSingleOptionAnswer(ans, len(q.Options)) {
				return false
			}
		case TypeMultiOptionQuestion:
			if !poll.isValidMultiOptionAnswer(ans, len(q.Options)) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func (poll *Poll) isValidSingleOptionAnswer(ans Answer, numOptions int) bool {
	ansSlice := []int(ans)
	if len(ansSlice) > 1 {
		return false
	}
	for _, option := range ansSlice {
		if option < 0 || option > numOptions {
			return false
		}
	}
	return true
}

func (poll *Poll) isValidMultiOptionAnswer(ans Answer, numOptions int) bool {
	ansSlice := []int(ans)
	if len(ansSlice) > numOptions {
		return false
	}
	for _, option := range ansSlice {
		if option < 0 || option >= numOptions {
			return false
		}
	}
	return true
}
