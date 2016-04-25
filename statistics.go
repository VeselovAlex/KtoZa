package KtoZa

// Statistics represents simple statistics for poll
// based on answer count
type Statistics struct {
	PollTitle string     `json:"polltitle,omitempty"`
	Answers   []statItem `json:"answers"`
}

type statItem struct {
	Text    string           `json:"text"`
	Options []statItemOption `json:"options"`
}

type statItemOption struct {
	Text  string `json:"text,omitempty"`
	Count int    `json:"count"`
}

// NewStatistics generates new statistics fo current poll
func NewStatistics(poll *Poll) *Statistics {
	res := new(Statistics)
	res.PollTitle = poll.Title
	res.Answers = make([]statItem, len(poll.Questions))

	for i, question := range poll.Questions {
		res.Answers[i] = newStatItem(question)
	}
	return res
}

func newStatItem(question pollQuestion) statItem {
	item := statItem{
		Text:    question.Text,
		Options: make([]statItemOption, len(question.Options)),
	}
	for i, opt := range question.Options {
		item.Options[i] = statItemOption{
			Text:  opt.Text,
			Count: 0,
		}
	}
	return item
}

// JoinWith adds results from stat to current statistics
func (s *Statistics) JoinWith(stat *Statistics) {
	switch {
	case s.PollTitle != stat.PollTitle:
		// Return
	case len(s.Answers) == 0:
		s.Answers = stat.Answers
	case len(s.Answers) != len(stat.Answers):
		// Return
	default:
		s.forceJoin(stat)
	}
}

// TODO Add error handling
func (s *Statistics) forceJoin(stat *Statistics) {
	for i, src := range stat.Answers {
		dest := s.Answers[i]
		for j, opt := range src.Options {
			dest.Options[j].Count += opt.Count
		}
	}
}

func (s *Statistics) ApplyAnswer(ans *Answer) {
	// TODO add verification
	answers := []answerOptions(*ans)
	if len(s.Answers) != len(answers) {
		return
	}
	for i, a := range answers {
		stat := s.Answers[i]
		opts := []string(a)
		if len(stat.Options) != len(opts) {
			continue
		}
		for _, opt := range opts {
			for j, statOpt := range stat.Options {
				if statOpt.Text == opt {
					stat.Options[j].Count++
				}
			}
		}
	}
}
