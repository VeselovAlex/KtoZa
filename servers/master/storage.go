package main

import (
	core "KtoZa"
)

type pollStorage struct {
	//Dummy
	poll *core.Poll
}

func (s pollStorage) Create(poll *core.Poll) error {
	s.poll = poll
	return nil
}

func (s pollStorage) Get() (*core.Poll, error) {
	return s.poll, nil
}

func (s pollStorage) Update(poll *core.Poll) error {
	s.poll = poll
	return nil
}

func (s pollStorage) Delete(poll *core.Poll) error {
	s.poll = nil
	return nil
}

var PollStorage = &pollStorage{}
