package distconf

import "sync"

type watchlist struct {
	watches []func()
	mu      sync.Mutex
}

func (w *watchlist) update() {
	for _, watch := range w.copyWatches() {
		watch()
	}
}

func (w *watchlist) copyWatches() []func() {
	w.mu.Lock()
	defer w.mu.Unlock()
	m := make([]func(), len(w.watches))
	copy(m, w.watches)
	return m
}

// Watch adds a watch for changes to this structure
func (w *watchlist) Watch(watch func()) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.watches = append(w.watches, watch)
}
