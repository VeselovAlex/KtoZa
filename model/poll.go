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

// Question представляет структуру вопроса
//
// Допустимы следующие типы вопроса (поле Type):
// * single-option    Вопрос с выбором одного варианта ответа
// * multi-option     Вопрос с выбором одного или нескольких вариантов ответа
//
type Question struct {
	Text    string   `json:"text"`
	Type    string   `json:"type"`
	Options []string `json:"options"`
}
