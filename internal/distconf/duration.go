package distconf

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"
)

type durationConf struct {
	Duration
	defaultVal time.Duration
	logger     Logger
}

// Duration is a duration type config inside a Config.
type Duration struct {
	watchlist
	currentVal int64
}

// Get the string in this config variable
func (s *Duration) Get() time.Duration {
	return time.Duration(atomic.LoadInt64(&s.currentVal))
}

func (s *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Get())
}

// Update the contents of Duration to the new value
func (s *durationConf) Update(newValue []byte) error {
	oldValue := s.Get()
	if newValue == nil {
		atomic.StoreInt64(&s.currentVal, int64(s.defaultVal))
	} else {
		newValDuration, err := time.ParseDuration(string(newValue))
		if err != nil {
			s.logger("", err, fmt.Sprintf("invalid duration string: %s", newValue))
			atomic.StoreInt64(&s.currentVal, int64(s.defaultVal))
		} else {
			atomic.StoreInt64(&s.currentVal, int64(newValDuration))
		}
	}
	if oldValue != s.Get() {
		s.update()
	}

	return nil
}

func (s *durationConf) GenericGet() interface{} {
	return s.Get()
}
