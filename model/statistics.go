// Александр Веселов <veselov143@gmail.com>
// СПбГУ, Математико-механический факультет, гр. 442
// Май, 2016 г.

// statistics.go содержит реализацию модели статистики системы KtoZa
package model

import (
	"errors"
	"time"
)

// Statistics представляет экземпляр статистики для опроса
type Statistics struct {
	// Время последнего обновления статистики
	LastUpdate time.Time `json:"date"`

	// Статистика по отдельным вопросам
	Questions []QuestionStat `json:"questions"`

	// Число принятых ответов
	RespondentsCount int `json:"respondents"`
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
	// Число ответов, содержащих данный вариант
	Count int `json:"count"`
}

// CreateStatisticsFor создает и инициализирует объект статистики
// в соответствии с заданным опросом
func CreateStatisticsFor(poll *Poll) *Statistics {
	if poll == nil {
		return nil
	}
	stat := &Statistics{
		LastUpdate: time.Now(),
		Questions:  make([]QuestionStat, len(poll.Questions)),
	}
	for i, question := range poll.Questions {
		stat.Questions[i].Options = make([]OptionStat, len(question.Options))
	}
	return stat
}

// JoinWith объединяет данную статистику с текущей, суммируя ответы. Если
// статистики были объединены, обновляет LastUpdate и возвращает true.
// Если текущая или заданная статистика == nil, или статистики объединить
// невозможно, возвращает false
func (stat *Statistics) JoinWith(other *Statistics) bool {
	if stat == nil || other == nil || !other.IsJoinableWith(stat) {
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

// JoinStatistics объединяет текущие статистики и возвращает
// новый экземпляр объединенной статистики. Функция не изменяет
// передаваемые статистики. Если хотя бы один из аргументов == nil
// или статистики невозможно объединить, возвращает nil
func JoinStatistics(one, other *Statistics) *Statistics {
	if one == nil || other == nil || !other.IsJoinableWith(one) {
		return nil
	}
	ret := &Statistics{}
	one.CopyTo(ret)
	ret.JoinWith(other)
	return ret
}

// CopyTo копирует текущий экземпляр статистики в dst.
// Если dst == nil, вызывает панику
func (stat *Statistics) CopyTo(dst *Statistics) {
	if dst == nil {
		panic(errors.New("Statistics copy dst is nil"))
	}
	dst.LastUpdate = stat.LastUpdate
	dst.RespondentsCount = stat.RespondentsCount

	dst.Questions = make([]QuestionStat, len(stat.Questions))
	copy(dst.Questions, stat.Questions)
	for i, q := range stat.Questions {
		dst.Questions[i].Options = make([]OptionStat, len(q.Options))
		copy(dst.Questions[i].Options, q.Options)
	}
}

func (qs *QuestionStat) joinWith(other QuestionStat) {
	for i, os := range other.Options {
		qs.Options[i].Count += os.Count
	}
	qs.AnswersCount += other.AnswersCount
}

// IsJoinableWith проверяет возможность объединения текущей статистики с
// заданной. Возвращает true, если текущая и заданная статистики не nil,
// число опросов и число вариантов ответов на каждый вопрос совпадают,
// иначе возвращает false
func (stat *Statistics) IsJoinableWith(other *Statistics) bool {
	if stat == nil || other == nil {
		return false
	}
	if len(other.Questions) != len(stat.Questions) {
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
	hasUpdates := len(sliceAns) > 0
	for i, ans := range sliceAns {
		stat.Questions[i].applyAnswer(ans)

	}

	if hasUpdates {
		stat.RespondentsCount++
		stat.LastUpdate = time.Now()
	}
	return hasUpdates
}

func (qs *QuestionStat) applyAnswer(ans Answer) {
	ansAsSlice := []int(ans)
	for _, option := range ansAsSlice {
		qs.Options[option].Count++
	}
	qs.AnswersCount++
}

func (stat *Statistics) isValidAnswerSet(ansSet AnswerSet) bool {
	ansSetAsSlice := []Answer(ansSet)
	if len(stat.Questions) != len(ansSetAsSlice) {
		return false
	}
	return true
}
