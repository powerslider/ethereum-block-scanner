package memory

import "sync"

// MultiMap represents a concurrent multimap data structure with one key and a slice of values.
type MultiMap[K comparable, V any] struct {
	sync.RWMutex
	m map[K][]V
}

// New instantiates a new multimap.
func New[K comparable, V any]() *MultiMap[K, V] {
	return &MultiMap[K, V]{m: make(map[K][]V)}
}

// Get searches the element in the multimap by key.
// It returns its value or nil if key is not found in multimap.
// Second return parameter is true if key was found, otherwise false.
func (m *MultiMap[K, V]) Get(key K) ([]V, bool) {
	m.RLock()
	values, found := m.m[key]
	m.RUnlock()

	return values, found
}

// Put stores a key-value pair in the multimap.
func (m *MultiMap[K, V]) Put(key K, value V) {
	m.Lock()
	m.m[key] = append(m.m[key], value)
	m.Unlock()
}

// PutAll stores a key-value pair in then multimap for each of the values, all using the same key.
func (m *MultiMap[K, V]) PutAll(key K, values []V) {
	m.Lock()
	for _, value := range values {
		m.Put(key, value)
	}
	m.RUnlock()
}
