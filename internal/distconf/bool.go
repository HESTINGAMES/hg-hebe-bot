package distconf

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)

type boolConf struct {
	Bool
	defaultVal int32
}

// Bool is a Boolean type config inside a Config.  It uses strconv.ParseBool to parse the conf
// contents as either true for false
type Bool struct {
	currentVal int32
	watchlist
}

// Update the contents of Bool to the new value
func (s *boolConf) Update(newValue []byte) error {
	oldValue := s.Get()
	if newValue == nil {
		atomic.StoreInt32(&s.currentVal, s.defaultVal)
	} else {
		newValueStr := string(newValue)
		if parsedBool, err := strconv.ParseBool(newValueStr); err != nil {
			atomic.StoreInt32(&s.currentVal, s.defaultVal)
		} else if parsedBool {
			atomic.StoreInt32(&s.currentVal, 1)
		} else {
			atomic.StoreInt32(&s.currentVal, 0)
		}
	}
	if oldValue != s.Get() {
		s.update()
	}
	return nil
}

// Get the boolean in this config variable
func (s *Bool) Get() bool {
	return atomic.LoadInt32(&s.currentVal) != 0
}

func (s *Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Get())
}

func (s *boolConf) GenericGet() interface{} {
	return s.Get()
}
