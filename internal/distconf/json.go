package distconf

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

// JSONConfig reads configuration from a JSON stream
type JSONConfig struct {
	vals    map[string][]byte
	watches map[string][]func(string)
	mu      sync.RWMutex
}

// RefreshFile reloads the configuration from a file
func (j *JSONConfig) RefreshFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	if err := j.Refresh(f); err != nil {
		return err
	}
	return f.Close()
}

// Refresh loads the configuration from a Reader
func (j *JSONConfig) Refresh(input io.Reader) error {
	fileContents := map[string]string{}
	if err := json.NewDecoder(input).Decode(&fileContents); err != nil {
		return err
	}
	j.mu.Lock()
	newVals := make(map[string][]byte, len(fileContents))
	for k, v := range fileContents {
		newVals[k] = []byte(v)
	}
	j.vals = newVals
	j.mu.Unlock()
	for k, cbs := range j.watches {
		for _, cb := range cbs {
			cb(k)
		}
	}
	return nil
}

// Get returns the key's value as read by JSON
func (j *JSONConfig) Get(key string) ([]byte, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	ret, exists := j.vals[key]
	if exists {
		return ret, nil
	}
	return nil, nil
}

// Close does nothing and exists to satisfy an interface
func (j *JSONConfig) Close() {
}

// Watch updates callback when a value changes inside the file
func (j *JSONConfig) Watch(key string, callback func(string)) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.watches == nil {
		j.watches = make(map[string][]func(string))
	}
	_, existing := j.watches[key]
	if !existing {
		j.watches[key] = []func(string){}
	}
	j.watches[key] = append(j.watches[key], callback)
	return nil
}
