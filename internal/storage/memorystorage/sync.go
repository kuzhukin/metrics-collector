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
