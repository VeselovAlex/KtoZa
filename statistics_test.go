package KtoZa

import (
	"encoding/json"
	"testing"
)

func TestStatisticsCreate(t *testing.T) {
	testStat := NewStatistics(testPoll)
	verifyEquals(t, "stat title", testPoll.Title, testStat.PollTitle)

	j, e := json.Marshal(testStat)
	t.Logf("JSON: %s\nError: %v\n", j, e)
}

func TestStatisticsJoin(t *testing.T) {
	poll := Poll{}
	poll.Questions = []pollQuestion{
		pollQuestion{
			Options: []pollOption{
				pollOption{},
			},
		},
	}

	stat := NewStatistics(poll)
	stat2 := NewStatistics(poll)
	stat2.Answers[0].Options[0].Count = 3
	stat.JoinWith(stat2)
	verifyEquals(t, "stat join", 3, stat.Answers[0].Options[0].Count)
}

func verifyEquals(t *testing.T, verifyOn string, exp, got interface{}) {
	if exp != got {
		t.Errorf("Error on verification of %s: expected %v, got %v\n",
			verifyOn,
			exp,
			got,
		)
	}
}
