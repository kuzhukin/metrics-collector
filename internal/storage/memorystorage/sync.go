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

func (s *syncMemoryStorage[T]) Update(k string, value T) {
	s.Lock()
	defer s.Unlock()

	s.storage[k] = value
}

func (s *syncMemoryStorage[T]) Get(k string) (T, bool) {
	s.Lock()
	defer s.Unlock()

	v, ok := s.storage[k]

	return v, ok
}
