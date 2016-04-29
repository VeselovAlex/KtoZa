package model

import "time"

// Statistics представляет экземпляр статистики для опроса
type Statistics struct {
	// Время последнего обновления статистики
	LastUpdate time.Time `json:"date"`

	// Статистика по отдельным вопросам
	Questions []QuestionStat `json:"questions"`

	// Число принятых ответов
	RespondentsCount int `json:"respondents"`

	poll *Poll
}

// QuestionStat представляет данные статистики отдельного вопроса
type QuestionStat struct {
	// Число принятых ответов на вопрос
	AnswersCount int `json:"answerCount"`

	// Статистика по вариантам ответа
	Options []OptionStat `json:"options"`
}

// OptionStat представляет статистику по отдельному варианту ответа
type OptionStat struct {
	Count int `json:"count"`
}

// CreateStatisticsFor создает и инициализирует объект статистики
// в соответствии с заданным опросом
func CreateStatisticsFor(poll *Poll) *Statistics {
	stat := &Statistics{
		LastUpdate: time.Now(),
		Questions:  make([]QuestionStat, len(poll.Questions)),
		poll:       poll,
	}
	for i, question := range poll.Questions {
		stat.Questions[i].Options = make([]OptionStat, len(question.Options))
	}
	return stat
}

// JoinWith объединяет данную статистику с текущей, суммируя ответы. Если
// статистики были объединены, обновляет LastUpdate и возвращает true
func (stat *Statistics) JoinWith(other *Statistics) bool {
	if !other.isJoinableWith(stat) {
		return false
	}

	hasUpdates := false

	for i, que := range other.Questions {
		if que.AnswersCount != 0 {
			stat.Questions[i].joinWith(que)
			hasUpdates = true
		}
	}

	if hasUpdates {
		stat.RespondentsCount += other.RespondentsCount
		stat.LastUpdate = time.Now()
	}
	return hasUpdates
}

func (qs *QuestionStat) joinWith(other QuestionStat) {
	for i, os := range other.Options {
		qs.Options[i].Count += os.Count
	}
	qs.AnswersCount += other.AnswersCount
}

func (stat *Statistics) isJoinableWith(other *Statistics) bool {
	if len(other.Questions) != len(stat.Questions) {
		// Статистика не соответствует текущей
		return false
	}

	for i, q := range stat.Questions {
		if len(q.Options) != len(other.Questions[i].Options) {
			return false
		}
	}

	return true
}

// ApplyAnswerSet проверяет корректность набора ответов и в случае
// корректности набора применяет его к текущей статистике. Если
// набор был применен, обновляет LastUpdate и возвращает true
func (stat *Statistics) ApplyAnswerSet(ans AnswerSet) bool {
	if !stat.isValidAnswerSet(ans) {
		return false
	}

	sliceAns := []Answer(ans)
	hasUpdates := false
	for i, ans := range sliceAns {
		stat.Questions[i].applyAnswer(ans)
		hasUpdates = true
	}

	if hasUpdates {
		stat.RespondentsCount++
		stat.LastUpdate = time.Now()
	}
	return hasUpdates
}

func (qs *QuestionStat) applyAnswer(ans Answer) {
	ansAsSlice := []int(ans)
	for i := range ansAsSlice {
		qs.Options[i].Count++
	}
	qs.AnswersCount++
}

func (stat *Statistics) isValidAnswerSet(ansSet AnswerSet) bool {
	ansSetAsSlice := []Answer(ansSet)
	if len(stat.Questions) != len(ansSetAsSlice) {
		return false
	}

	for range ansSetAsSlice {
		//checkValidAnswer
	}
	return true
}
