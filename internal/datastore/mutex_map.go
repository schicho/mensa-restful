package datastore

import "sync"

type mutexMap[K comparable, V any] struct {
	sync.RWMutex
	internal map[K]V
}

func newMutexMap[K comparable, V any]() *mutexMap[K,V] {
	return &mutexMap[K,V]{
		internal: make(map[K]V),
	}
}

func (rm *mutexMap[K,V]) Load(key K) (V, bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}

func (rm *mutexMap[K,V]) Delete(key K) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *mutexMap[K,V]) Store(key K, value V) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}