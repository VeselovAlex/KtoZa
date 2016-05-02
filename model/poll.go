package model

import "time"

// Poll представляет структуру данных опроса системы KtoZa
type Poll struct {
	Title     string       `json:"title"`
	Caption   string       `json:"caption"`
	Events    EventTimings `json:"events"`
	Questions []Question   `json:"questions"`
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
