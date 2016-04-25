package main

import (
	core "KtoZa"
)

type statController struct {
	data *core.Statistics
}

func (s *statController) Init(stats *core.Statistics) {
	s.data = stats
}

func (s statController) Get() *core.Statistics {
	return s.data
}

func (s statController) JoinWith(stats *core.Statistics) {
	s.data.JoinWith(stats)
}

var StatisticsCtrl statController
