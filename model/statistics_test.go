package model

import (
	"testing"
	"time"
)

func TestJoin(t *testing.T) {
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
	stat1.Join(stat2)

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
