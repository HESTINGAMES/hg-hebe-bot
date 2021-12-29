package distconf

import (
	"bytes"
	"sync"
)

// ReaderCache is a type of distconf reader that simply caches previous Get() values and updates watchers when a value
// changes
type ReaderCache struct {
	lastKnownValues map[string][]byte
	watchedKeys     map[string][]func(string)
	mu              sync.Mutex
}

// NotifyWatchers is called by users of ReaderCache to notify registered Watches that a value changed
func (s *ReaderCache) notifyWatchers(key string, newVal []byte) {
	s.mu.Lock()
	lastKnown := s.lastKnownValues[key]
	if s.lastKnownValues == nil {
		s.lastKnownValues = make(map[string][]byte)
	}
	s.lastKnownValues[key] = newVal
	toNotify := make([]func(string), 0, len(s.watchedKeys[key]))
	toNotify = append(toNotify, s.watchedKeys[key]...)

	s.mu.Unlock()
	if !bytes.Equal(lastKnown, newVal) {
		for _, notifier := range toNotify {
			notifier(key)
		}
	}
}

// Watch registers a distconf watcher
func (s *ReaderCache) Watch(key string, callback func(string)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.watchedKeys == nil {
		s.watchedKeys = make(map[string][]func(string))
	}
	s.watchedKeys[key] = append(s.watchedKeys[key], callback)
	return nil
}

// CopyWatchedKeys copies the internal watched keys array
func (s *ReaderCache) copyWatchedKeys() map[string][]func(string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.watchedKeys == nil {
		return nil
	}
	ret := make(map[string][]func(string), len(s.watchedKeys))
	for k, v := range s.watchedKeys {
		ret[k] = make([]func(string), 0, len(v))
		ret[k] = append(ret[k], v...)
	}
	return ret
}

// ReaderCacheNotify notifies a readercache that a value changed.  This function isn't public on the ReaderCache
// so structs can embed a ReaderCache directly without exposing notifyWatchers.
func ReaderCacheNotify(s *ReaderCache, key string, newVal []byte) {
	s.notifyWatchers(key, newVal)
}

// ReaderCacheRefresh calls a Get on every value watched by the ReaderCache
func ReaderCacheRefresh(s *ReaderCache, r Reader, OnFetchErr func(err error, key string)) {
	keys := s.copyWatchedKeys()
	for keyName := range keys {
		_, err := r.Get(keyName)
		if err != nil && OnFetchErr != nil {
			OnFetchErr(err, keyName)
		}
	}
}
