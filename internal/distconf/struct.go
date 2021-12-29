package distconf

import (
	"encoding/json"
	"reflect"
	"sync/atomic"
)

type structConf struct {
	Struct
	defaultVal interface{}
	t          reflect.Type
}

// Struct allows distconf to atomically update
// sets of variables. The structure provided must be represented as a
// json object in the backend. The DefaultValue provided sets the type that
// the backend josn will be unmarshelled into.
type Struct struct {
	watchlist
	currentVal atomic.Value
}

// Update the contents of Struct to the new value
func (s *structConf) Update(newValue []byte) error {
	oldValue := s.currentVal.Load()
	if newValue == nil {
		s.currentVal.Store(s.defaultVal)
	} else {
		val := reflect.New(s.t)
		addrVal := val.Interface()
		err := json.Unmarshal(newValue, addrVal)
		if err != nil {
			s.currentVal.Store(s.defaultVal)
		} else {
			s.currentVal.Store(reflect.Indirect(val).Interface())
		}
	}
	if reflect.DeepEqual(oldValue, s.currentVal.Load()) {
		s.update()
	}
	return nil
}

// Get the current value. For structs this *must* be cast
// to the same type as the DefaultValue provided. Not matching
// the type will result in a runtime panic.
func (s *Struct) Get() interface{} {
	return s.currentVal.Load()
}

func (s *Struct) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Get())
}

func (s *structConf) GenericGet() interface{} {
	return s.currentVal.Load()
}
