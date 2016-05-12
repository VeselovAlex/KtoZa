// Александр Веселов <veselov143@gmail.com>
// СПбГУ, Математико-механический факультет, гр. 442
// Май, 2016 г.

// statistics_test.go содержит unit-тесты модели статистики системы KtoZa
package model

import (
	"testing"
	"time"
)

func TestJoinWith(t *testing.T) {
	now := time.Now()
	stat1 := &Statistics{
		LastUpdate:       now,
		RespondentsCount: 10,
		Questions: []QuestionStat{
			QuestionStat{
				AnswersCount: 5,
				Options: []OptionStat{
					OptionStat{
						Count: 5,
					},
				},
			},
		},
	}

	stat2 := &Statistics{
		LastUpdate:       time.Now(),
		RespondentsCount: 10,
		Questions: []QuestionStat{
			QuestionStat{
				AnswersCount: 5,
				Options: []OptionStat{
					OptionStat{
						Count: 5,
					},
				},
			},
		},
	}

	// Чтобы проверить изменение времени последнего обновления
	time.Sleep(100 * time.Millisecond)
	stat1.JoinWith(stat2)

	// Все счетчики должны удвоиться
	if stat1.RespondentsCount != 2*stat2.RespondentsCount {
		t.Errorf("RespondentsCount mismatch: expected %d, got %d\n",
			stat2.RespondentsCount,
			stat1.RespondentsCount)
	}

	// Все счетчики должны удвоиться
	for i, q := range stat1.Questions {
		if q.AnswersCount != 2*stat2.Questions[i].AnswersCount {
			t.Errorf("AnswerCount mismatch: expected %d, got %d\n",
				stat2.Questions[i].AnswersCount,
				q.AnswersCount)
		}

		for j, o := range q.Options {
			if q.AnswersCount != 2*stat2.Questions[i].Options[j].Count {
				t.Errorf("Option Count mismatch: expected %d, got %d\n",
					stat2.Questions[i].Options[j].Count,
					o.Count)
			}
		}
	}

	if !stat1.LastUpdate.After(now) {
		t.Errorf("LastUpdate has not been updated")
	}
}

func TestApplyAnswerSet(t *testing.T) {
	now := time.Now()
	poll := &Poll{
		Questions: []Question{
			Question{
				Type:    TypeSingleOptionQuestion,
				Options: []string{"0", "1"},
			},
		},
	}

	stat := CreateStatisticsFor(poll)

	ansSet := AnswerSet([]Answer{
		Answer([]int{1}),
	})

	// Чтобы проверить изменение времени последнего обновления
	time.Sleep(100 * time.Millisecond)
	stat.ApplyAnswerSet(ansSet)

	if stat.RespondentsCount != 1 {
		t.Errorf("RespondentsCount mismatch: expected %d, got %d\n",
			1,
			stat.RespondentsCount)
	}

	if stat.Questions[0].AnswersCount != 1 {
		t.Errorf("Answers count mismatch: expected %d, got %d\n",
			1,
			stat.Questions[0].AnswersCount)
	}

	if stat.Questions[0].Options[1].Count != 1 {
		t.Errorf("Option count mismatch: expected %d, got %d\n",
			1,
			stat.Questions[0].Options[1].Count)
	}

	if !stat.LastUpdate.After(now) {
		t.Errorf("LastUpdate has not been updated")
	}
}
