package distconf

import (
	"encoding/json"
	"sync/atomic"
)

type strConf struct {
	Str
	defaultVal string
}

// Str is a string type config inside a Config.
type Str struct {
	watchlist
	currentVal atomic.Value
}

// Update the contents of Str to the new value
func (s *strConf) Update(newValue []byte) error {
	oldValue := s.currentVal.Load().(string)
	if newValue == nil {
		s.currentVal.Store(s.defaultVal)
	} else {
		s.currentVal.Store(string(newValue))
	}
	if oldValue != s.Get() {
		s.update()
	}
	return nil
}

func (s *Str) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Get())
}

// Get the string in this config variable
func (s *Str) Get() string {
	return s.currentVal.Load().(string)
}

func (s *strConf) GenericGet() interface{} {
	return s.Get()
}
