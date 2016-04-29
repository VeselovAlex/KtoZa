package model

import (
	"testing"
	"time"
)

<<<<<<< HEAD
func TestJoinWith(t *testing.T) {
=======
func TestJoin(t *testing.T) {
>>>>>>> 053f099... statistics remodelled
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
<<<<<<< HEAD
	stat1.JoinWith(stat2)
=======
	stat1.Join(stat2)
>>>>>>> 053f099... statistics remodelled

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
<<<<<<< HEAD

func TestApplyAnswerSet(t *testing.T) {
	now := time.Now()
	stat := &Statistics{
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

	ansSet := AnswerSet([]Answer{
		Answer([]int{0}),
	})

	// Чтобы проверить изменение времени последнего обновления
	time.Sleep(100 * time.Millisecond)
	stat.ApplyAnswerSet(ansSet)

	if stat.RespondentsCount != 11 {
		t.Errorf("RespondentsCount mismatch: expected %d, got %d\n",
			11,
			stat.RespondentsCount)
	}

	if stat.Questions[0].AnswersCount != 6 {
		t.Errorf("Answers count mismatch: expected %d, got %d\n",
			6,
			stat.Questions[0].AnswersCount)
	}

	if stat.Questions[0].Options[0].Count != 6 {
		t.Errorf("Option count mismatch: expected %d, got %d\n",
			6,
			stat.Questions[0].Options[0].Count)
	}

	if !stat.LastUpdate.After(now) {
		t.Errorf("LastUpdate has not been updated")
	}
}
=======
>>>>>>> 053f099... statistics remodelled
