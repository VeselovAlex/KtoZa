package model

import "time"

type Statistics struct {
	LastUpdate       time.Time      `json:"date"`
	Questions        []QuestionStat `json:"questions"`
	RespondentsCount int            `json:"respondents"`
}

type QuestionStat struct {
	AnswersCount int          `json:"answerCount"`
	Options      []OptionStat `json:"options"`
}

type OptionStat struct {
	Count int `json:"count"`
}

// Join объединяет данную статистику с текущей, суммируя ответы.
// Возвращает true, если статистики были объединены
func (stat *Statistics) Join(other *Statistics) bool {
	if len(other.Questions) != len(stat.Questions) {
		// Статистика не соответствует текущей
		return false
	}

	hasUpdates := false

	for i, que := range other.Questions {
		// Соответствуюий вопрос из текущей статистики
		matchQ := stat.Questions[i]
		if len(que.Options) != len(matchQ.Options) {
			// Вопрос не соответствует текущему
			continue
		}

		if que.AnswersCount != 0 {
			// Вопрос обновлен
			stat.Questions[i].AnswersCount += que.AnswersCount
			for j, opt := range que.Options {
				// Суммируем ответы
				matchQ.Options[j].Count += opt.Count
			}
			hasUpdates = true
		}

	}

	if hasUpdates {
		stat.RespondentsCount += other.RespondentsCount
		stat.LastUpdate = time.Now()
	}
	return hasUpdates
}
