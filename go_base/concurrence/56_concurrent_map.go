package concurrence

import "sync"

type ConcurrentMap[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func NewConcurrentMap[K comparable, V any](cap int) *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		data: make(map[K]V, cap),
	}
}

func (cm *ConcurrentMap[K, V]) Store(key K, value V) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.data[key] = value
}

func (cm *ConcurrentMap[K, V]) Load(key K) (value V, exists bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	value, exists = cm.data[key]
	return
}

// 遍历map，根据回调函数f来处理<key, value>
func (cm *ConcurrentMap[K, V]) Range(f func(key K, value V) error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	for k, v := range cm.data {
		if err := f(k, v); err != nil {
			break
		}
	}
}
