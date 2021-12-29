package distconf

import (
	"sync"
)

// InMemory stores configuration in memory
type InMemory struct {
	vals    map[string][]byte
	watches map[string][]func(string)
	mu      sync.Mutex
}

// Get returns the stored value
func (m *InMemory) Get(key string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, exists := m.vals[key]
	if !exists {
		return nil, nil
	}
	return b, nil
}

// ListConfig returns a copy of the currently stored config values
func (m *InMemory) ListConfig() map[string][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	ret := make(map[string][]byte, len(m.vals))
	for k, v := range m.vals {
		ret[k] = v
	}
	return ret
}

// StoreConfig stores into memory the given values, only if there currently is not a copy
func (m *InMemory) StoreConfig(toStore map[string][]byte) {
	if m.vals == nil {
		m.vals = make(map[string][]byte, len(toStore))
	}
	for k, v := range toStore {
		m.mu.Lock()
		_, exists := m.vals[k]
		m.mu.Unlock()
		if !exists {
			m.write(k, v)
		}
	}
}

func (m *InMemory) write(key string, value []byte) {
	m.mu.Lock()
	if m.vals == nil {
		m.vals = make(map[string][]byte)
	}
	if value == nil {
		delete(m.vals, key)
	} else {
		m.vals[key] = value
	}
	m.mu.Unlock()

	for _, calls := range m.watches[key] {
		calls(key)
	}
}

// Write updates the in memory value
func (m *InMemory) Write(key string, value []byte) error {
	m.write(key, value)
	return nil
}

// Watch calls callback when a write happens on the key
func (m *InMemory) Watch(key string, callback func(string)) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.watches == nil {
		m.watches = make(map[string][]func(string))
	}
	_, existing := m.watches[key]
	if !existing {
		m.watches[key] = []func(string){callback}
		return nil
	}
	m.watches[key] = append(m.watches[key], callback)
	return nil
}

// Close does nothing and exists to satisfy an interface
func (m *InMemory) Close() {
}
