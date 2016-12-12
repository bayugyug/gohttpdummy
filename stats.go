package main

import (
	"sync"
)

//StatsHelper - are holder of summary
type StatsHelper struct {
	locker sync.Mutex
	stats  map[string]int
}

//StatsHelperNew - return new object
func StatsHelperNew() (n *StatsHelper) {
	//init
	return &StatsHelper{stats: make(map[string]int)}
}

func (s *StatsHelper) getStatsList() map[string]int {
	s.locker.Lock()
	defer s.locker.Unlock()
	statz := s.stats
	return statz
}

func (s *StatsHelper) setStats(prefix string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	if prefix == "" {
		return
	}
	s.stats[prefix]++
}

func (s *StatsHelper) getStats(prefix string) int {
	s.locker.Lock()
	defer s.locker.Unlock()
	if prefix == "" {
		return 0
	}
	m := s.stats[prefix]
	return m
}
