package memorystorage

import "sync"

func newSyncMemoryStorage[T int64 | float64]() syncMemoryStorage[T] {
	return syncMemoryStorage[T]{
		storage: make(map[string]T),
	}
}

type syncMemoryStorage[T int64 | float64] struct {
	sync.Mutex
	storage map[string]T
}

func (s *syncMemoryStorage[T]) Write(k string, value T) {
	s.Lock()
	defer s.Unlock()

	s.storage[k] = value
}

func (s *syncMemoryStorage[T]) Sum(k string, value T) {
	s.Lock()
	defer s.Unlock()

	s.storage[k] += value
}

func (s *syncMemoryStorage[T]) Get(k string) (T, bool) {
	s.Lock()
	defer s.Unlock()

	v, ok := s.storage[k]

	return v, ok
}

func (s *syncMemoryStorage[T]) GetAll() map[string]T {
	s.Lock()
	defer s.Unlock()

	copy := make(map[string]T, len(s.storage))

	for k, v := range s.storage {
		copy[k] = v
	}

	return copy
}
