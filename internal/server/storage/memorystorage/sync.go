package memorystorage

import "sync"

func NewSyncStorage[T int64 | float64]() SyncStorage[T] {
	return SyncStorage[T]{
		storage: make(map[string]T),
	}
}

type SyncStorage[T int64 | float64] struct {
	sync.RWMutex
	storage map[string]T
}

func (s *SyncStorage[T]) Write(k string, value T) {
	s.Lock()
	defer s.Unlock()

	s.storage[k] = value
}

func (s *SyncStorage[T]) Sum(k string, value T) {
	s.Lock()
	defer s.Unlock()

	s.storage[k] += value
}

func (s *SyncStorage[T]) Get(k string) (T, bool) {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.storage[k]

	return v, ok
}

func (s *SyncStorage[T]) GetAll() map[string]T {
	s.RLock()
	defer s.RUnlock()

	copy := make(map[string]T, len(s.storage))

	for k, v := range s.storage {
		copy[k] = v
	}

	return copy
}

func (s *SyncStorage[T]) SetAll(m map[string]T) {
	s.Lock()
	defer s.Unlock()

	s.storage = m
}
