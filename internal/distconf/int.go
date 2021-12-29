package distconf

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)

type intConf struct {
	Int
	defaultVal int64
}

// Int is an integer type config inside a Config.
type Int struct {
	watchlist
	currentVal int64
}

// Get the integer in this config variable
func (c *Int) Get() int64 {
	return atomic.LoadInt64(&c.currentVal)
}

func (c *Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Get())
}

// Update the content of this config variable to newValue.
func (c *intConf) Update(newValue []byte) error {
	oldValue := c.Get()
	if newValue == nil {
		atomic.StoreInt64(&c.currentVal, c.defaultVal)
	} else {
		newValueInt, err := strconv.ParseInt(string(newValue), 10, 64)
		if err != nil {
			return err
		}
		atomic.StoreInt64(&c.currentVal, newValueInt)
	}
	if oldValue != c.Get() {
		c.update()
	}

	return nil
}

func (c *intConf) GenericGet() interface{} {
	return c.Get()
}
