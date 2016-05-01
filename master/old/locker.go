package main

import "sync"

var StatisticsLocker = &locker{}
var PollLocker = &locker{}

type locker struct {
	lock sync.RWMutex
}

func (l *locker) Do(f func()) {
	l.lock.Lock()
	defer l.lock.Unlock()
	f()
}

func (l *locker) DoReadOnly(f func()) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	f()
}
