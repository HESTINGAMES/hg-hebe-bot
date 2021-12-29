package distconf

import (
	"bytes"
	"encoding/json"
	"io"
)

// DynamicReader is a reader that can also change
type DynamicReader interface {
	Reader
	Dynamic
}

// CachedReader allows bulk loading values and caching the result
type CachedReader struct {
	cache    InMemory
	Fallback DynamicReader
}

// Get will try fetching it from source.  If you can, update in memory.  If you can't, try to find it in memory and return that
// instead.
func (m *CachedReader) Get(key string) ([]byte, error) {
	ret, err := m.Fallback.Get(key)
	if err == nil {
		err = m.cache.Write(key, ret)
		return ret, err
	}
	ret, err2 := m.cache.Get(key)
	if err2 != nil && ret != nil {
		return ret, nil
	}
	return nil, err
}

// Close ends the fallback endpoint
func (m *CachedReader) Close() {
	m.Fallback.Close()
}

// Watch forwards to the fallback
func (m *CachedReader) Watch(key string, callback func(string)) error {
	return m.Fallback.Watch(key, callback)
}

// ListConfig returns all the cached values
func (m *CachedReader) ListConfig() map[string][]byte {
	return m.cache.ListConfig()
}

// StoreConfig overwrites the stored config
func (m *CachedReader) StoreConfig(toStore map[string][]byte) {
	m.cache.StoreConfig(toStore)
}

// WriteTo updates the cache
func (m *CachedReader) WriteTo(w io.Writer) (int64, error) {
	toSave := m.cache.ListConfig()
	asJSON := make(map[string]string, len(toSave))
	for k, v := range toSave {
		asJSON[k] = string(v)
	}
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&asJSON); err != nil {
		return 0, err
	}
	n, err := w.Write(buf.Bytes())
	return int64(n), err
}

// ReadFrom loads a stored cache
func (m *CachedReader) ReadFrom(r io.Reader) (n int64, err error) {
	asJSON := make(map[string][]byte)
	buf := bytes.Buffer{}
	amnt, err := io.Copy(&buf, r)
	if err != nil {
		return amnt, err
	}

	if err := json.NewDecoder(&buf).Decode(&asJSON); err != nil {
		return 0, err
	}
	m.StoreConfig(asJSON)
	return amnt, nil
}
