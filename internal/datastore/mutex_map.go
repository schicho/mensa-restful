package datastore

import "sync"

type mutexMap struct {
	sync.RWMutex
	internal map[int]cacheddata
}

func newMutexMap() *mutexMap {
	return &mutexMap{
		internal: make(map[int]cacheddata),
	}
}

func (rm *mutexMap) Load(key int) (cacheddata, bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}

func (rm *mutexMap) Delete(key int) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *mutexMap) Store(key int, value cacheddata) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}