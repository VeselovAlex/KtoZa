package KtoZa

import (
	"encoding/json"
	"testing"
	"time"
)

var testPoll = Poll{
	pollHeader{
		Title:   "Hello",
		Caption: "Hello test",
		StartAt: time.Now(),
	},
	[]pollQuestion{
		pollQuestion{
			Text: "Question 1",
			Options: []pollOption{
				pollOption{
					Text: "Yes",
				},
				pollOption{
					Text: "No",
				},
			},
		},
	},
}

func TestPollCreate(t *testing.T) {
	verifyEquals(t, "poll title", "Hello", testPoll.Title)
	verifyEquals(t, "poll question count", 1, len(testPoll.Questions))

	q := testPoll.Questions[0]
	verifyEquals(t, "poll question text", "Question 1", q.Text)
	verifyEquals(t, "poll question options", 2, len(q.Options))

	j, e := json.Marshal(testPoll)
	t.Logf("JSON: %s\nError: %v\n", j, e)
}
