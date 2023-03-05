package domain

import "sync"

type SagaSet struct {
	mu    sync.RWMutex
	sagas map[int64]*Saga
}

func NewSagaSet() *SagaSet {
	return &SagaSet{
		mu:    sync.RWMutex{},
		sagas: make(map[int64]*Saga),
	}
}

func (s *SagaSet) Get(id int64) *Saga {
	s.mu.Lock()
	defer s.mu.Unlock()

	saga, ok := s.sagas[id]
	if !ok {
		return nil
	}

	return saga
}

func (s *SagaSet) Register(id int64, saga *Saga) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.sagas[id]
	if !ok {
		s.sagas[id] = saga
	}
}

func (s *SagaSet) Remove(id int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.sagas[id]
	if ok {
		delete(s.sagas, id)
	}
}
