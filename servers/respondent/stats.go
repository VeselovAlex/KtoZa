package main

import (
	"net/http"
	"sync"

	core "KtoZa"
)

type statisticsCache struct {
	Answers chan *core.Answer

	nslock     sync.Mutex
	needSubmit bool

	statLock    sync.RWMutex
	curStat     *core.Statistics
	answerCache *core.Statistics
}

func (s *statisticsCache) Run() {
	go s.loopSubmit()
	for ans := range s.Answers {
		s.applyAnswer(ans)
	}
}

func (s *statisticsCache) applyAnswer(ans *core.Answer) {
	s.statLock.RLock()
	defer s.statLock.RUnlock()
	s.answerCache.ApplyAnswer(ans)
	s.curStat.ApplyAnswer(ans)
	s.setNeedSubmit()
}

func (s *statisticsCache) loopSubmit() {
	for {
		var ns bool
		s.nslock.Lock()
		ns = s.needSubmit
		s.needSubmit = false
		s.nslock.Unlock()

		if ns {
			newCache := core.NewStatistics(PollController.Get())
			s.statLock.Lock()
			cache := s.answerCache
			s.answerCache = newCache
			s.statLock.Unlock()
			err := Master.SendSubmit(cache)
			if err != nil {
				// Bad connection
				s.statLock.Lock()
				s.answerCache.JoinWith(cache)
				s.statLock.Unlock()
				s.setNeedSubmit()
			}
		}
	}
}

func (s *statisticsCache) setNeedSubmit() {
	s.nslock.Lock()
	s.needSubmit = true
	s.nslock.Unlock()
}

var StatisticsCache = statisticsCache{
	Answers: make(chan *core.Answer, 16),
}

func submitUserAnswer(w http.ResponseWriter, r *http.Request) {
	if !inTime() {
		http.NotFound(w, r)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed. Use POST instead", http.StatusMethodNotAllowed)
		return
	}
	a := &core.Answer{}
	a.FromRequest(r)
	StatisticsCache.Answers <- a
	w.WriteHeader(http.StatusAccepted)
}
