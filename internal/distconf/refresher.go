package distconf

import (
	"sync"
	"time"
)

// Refreshable is any configuration that needs to be periodically refreshed
type Refreshable interface {
	Refresh()
}

// ComboRefresher can refresh from multiple sources
type ComboRefresher []Refreshable

// Refresh calls all refreshes at once
func (c ComboRefresher) Refresh() {
	wg := sync.WaitGroup{}
	wg.Add(len(c))
	for _, r := range c {
		go func(r Refreshable) {
			defer wg.Done()
			r.Refresh()
		}(r)
	}
	wg.Wait()
}

// A Refresher can refresh the values inside any Refreshable object
type Refresher struct {
	WaitTime           *Duration
	ToRefresh          Refreshable
	startShouldEndChan chan struct{}
	startHasEndedChan  chan struct{}
	once               sync.Once
}

func (r *Refresher) setup() {
	r.once.Do(func() {
		r.startShouldEndChan = make(chan struct{})
		r.startHasEndedChan = make(chan struct{})
	})
}

// Setup inits the refresher
func (r *Refresher) Setup() error {
	r.setup()
	return nil
}

// Close ends the Refresher
func (r *Refresher) Close() error {
	r.setup()
	close(r.startShouldEndChan)
	return nil
}

// Done returns a channel that blocks until Close() is called or Start() is finished
func (r *Refresher) Done() <-chan struct{} {
	r.setup()
	return r.startHasEndedChan
}

// Start executing the refresher, waiting waitTime between ending, and a new start of the Refresh call
func (r *Refresher) Start() error {
	r.setup()
	defer close(r.startHasEndedChan)
	for {
		select {
		case <-r.startShouldEndChan:
			return nil
		case <-time.After(r.WaitTime.Get()):
			r.ToRefresh.Refresh()
		}
	}
}
