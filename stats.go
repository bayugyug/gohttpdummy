package main

import (
	"sync"
)

//StatsHelper - are holder of summary
type StatsHelper struct {
	locker sync.RWMutex
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
	if s.stats != nil {
		return s.stats
	}
	return nil
}

func (s *StatsHelper) setStats(prefix string) {
	if prefix == "" {
		return
	}
	s.locker.Lock()
	s.stats[prefix]++
	s.locker.Unlock()
}

func (s *StatsHelper) getStats(prefix string) int {
	if prefix == "" {
		return 0
	}
	s.locker.Lock()
	m := s.stats[prefix]
	s.locker.Unlock()
	return m
}
