package memorystorage

import (
	"fmt"
	"sync"
	"testing"
)

const testCount = 1000

func BenchmarkSyncStorageWrite(b *testing.B) {
	syncStorage := NewSyncStorage[int64]()
	values := prepareValues()
	syncMap := sync.Map{}

	b.ResetTimer()

	b.Run("sync-storage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for k, v := range values {
				syncStorage.Write(k, v)
			}
		}
	})
	b.Run("sync-map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for k, v := range values {
				syncMap.Store(k, v)
			}
		}
	})
}

func BenchmarkStoageWriteParallel(b *testing.B) {
	syncStorage := NewSyncStorage[int64]()
	values := prepareValues()

	b.ResetTimer()
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			for k, v := range values {
				syncStorage.Write(k, v)
			}
		}
	})
}

func BenchmarkSyncMapWriteParallel(b *testing.B) {
	values := prepareValues()

	syncMap := sync.Map{}

	b.ResetTimer()
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			for k, v := range values {
				syncMap.Store(k, v)
			}
		}
	})
}

func prepareValues() map[string]int64 {
	m := make(map[string]int64)

	for i := 0; i < testCount; i++ {
		m[fmt.Sprintf("metric-%d", i)] = int64(i + 1)
	}

	return m
}
