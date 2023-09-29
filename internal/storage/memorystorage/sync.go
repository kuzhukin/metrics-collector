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

type syncMemoryStorageUpdater[T int64 | float64] func(storage map[string]T)

func (s *syncMemoryStorage[T]) Update(fn syncMemoryStorageUpdater[T]) {
	s.Lock()
	defer s.Unlock()

	fn(s.storage)
}
