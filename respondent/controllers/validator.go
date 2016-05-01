package controllers

import "github.com/VeselovAlex/KtoZa/model"

type AnswerValidator interface {
	IsValid(model.Answer) bool
}

// NewAnswerValidatorFor создает валидатор для заданного вопроса
func NewAnswerValidatorFor(q *model.Question) AnswerValidator {
	switch q.Type {
	case model.TypeSingleOptionQuestion:
		return singleOptionValidator{max: len(q.Options) - 1}
	case model.TypeMultiOptionQuestion:
		return multiOptionValidator{numOptions: len(q.Options)}
	}
	return nil
}

type singleOptionValidator struct {
	max int
}

func (validator singleOptionValidator) IsValid(ans model.Answer) bool {
	ansSlice := []int(ans)
	if len(ansSlice) > 1 {
		// Больше одного варианта ответа
		return false
	}
	for _, option := range ansSlice {
		if option < 0 || option > validator.max {
			return false
		}
	}
	return true
}

type multiOptionValidator struct {
	numOptions int
}

func (validator multiOptionValidator) IsValid(ans model.Answer) bool {
	ansSlice := []int(ans)
	if len(ansSlice) > validator.numOptions {
		// Слишком много вариантов ответов
		return false
	}
	for _, option := range ansSlice {
		if option < 0 || option >= validator.numOptions {
			return false
		}
	}
	return true
}

type Validator struct {
	ansValidators []AnswerValidator
}

func (v Validator) IsValid(set model.AnswerSet) bool {
	asSlice := []model.Answer(set)
	// Проверка соответствия числа ответов числу вопросов
	if len(asSlice) != len(v.ansValidators) {
		return false
	}
	for i, ans := range asSlice {
		if !v.ansValidators[i].IsValid(ans) {
			return false
		}
	}
	return true
}

// NewValidatorFor создает валидатор для заданного опроса
func NewValidatorFor(poll *model.Poll) Validator {
	var validator Validator
	return validator
}
