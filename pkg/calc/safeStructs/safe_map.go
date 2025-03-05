package safeStructs

import (
	"sync"
)

type SafeMap struct {
	m   map[int]Expressions
	mux sync.RWMutex
}

type Expressions struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
	Result string `json:"result"`
}

func NewSafeMap() *SafeMap {
	return &SafeMap{m: make(map[int]Expressions), mux: sync.RWMutex{}}
}

func (s *SafeMap) Get(key int) Expressions {
	s.mux.RLock()
	defer s.mux.RUnlock()
	res, ok := s.m[key]
	if ok {
		return res
	} else {
		return Expressions{}
	}
}

func (s *SafeMap) In(key int) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	_, ok := s.m[key]
	if ok {
		return true
	} else {
		return false
	}
}

func (s *SafeMap) Set(key int, value Expressions) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m[key] = value
}

func (s *SafeMap) GetAll() []Expressions {
	s.mux.RLock()
	defer s.mux.RUnlock()
	var res []Expressions
	for _, v := range s.m {
		res = append(res, v)
	}
	return res
}
