package distconf

import (
	"encoding/json"
	"math"
	"strconv"
	"sync/atomic"
)

type floatConf struct {
	Float
	defaultVal float64
}

// Float is an float type config inside a Config.
type Float struct {
	watchlist
	// store as uint64 and convert on way in and out for atomicity
	currentVal uint64
}

// Get the float in this config variable
func (c *Float) Get() float64 {
	return math.Float64frombits(atomic.LoadUint64(&c.currentVal))
}

func (c *Float) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Get())
}

// Update the content of this config variable to newValue.
func (c *floatConf) Update(newValue []byte) error {
	oldValue := c.Get()
	if newValue == nil {
		atomic.StoreUint64(&c.currentVal, math.Float64bits(c.defaultVal))
	} else {
		newValueFloat, err := strconv.ParseFloat(string(newValue), 64)
		if err != nil {
			return err
		}
		atomic.StoreUint64(&c.currentVal, math.Float64bits(newValueFloat))
	}
	if oldValue != c.Get() {
		c.update()
	}

	return nil
}

func (c *floatConf) GenericGet() interface{} {
	return c.Get()
}
