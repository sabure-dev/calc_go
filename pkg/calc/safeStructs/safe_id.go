package safeStructs

import "sync"

type SafeId struct {
	Id  int
	mux sync.Mutex
}

func NewSafeId() *SafeId {
	return &SafeId{
		Id:  0,
		mux: sync.Mutex{},
	}
}

func (id *SafeId) Get() int {
	id.mux.Lock()
	defer id.mux.Unlock()
	id.Id++
	return id.Id
}
