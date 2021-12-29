package distconf

import (
	"expvar"
	"math"
	"reflect"
	"sync"
	"time"
)

// Logger is used to log debug information during distconf operation
type Logger func(key string, err error, msg string)

// Distconf gets configuration data from the first backing that has it
type Distconf struct {
	Logger  Logger
	Readers []Reader

	varsMutex      sync.Mutex
	registeredVars map[string]*registeredVariableTracker
}

type registeredVariableTracker struct {
	distvar        configVariable
	hasInitialized sync.Once
}

type configVariable interface {
	Update(newValue []byte) error
	// Get but on an interface return.  Oh how I miss you templates.
	GenericGet() interface{}
}

func (c *Distconf) logger(key string, err error, msg string) {
	if c.Logger != nil {
		c.Logger(key, err, msg)
	}
}

// Keys returns an array of all keys watched by distconf
func (c *Distconf) Keys() []string {
	c.varsMutex.Lock()
	defer c.varsMutex.Unlock()
	ret := make([]string, 0, len(c.registeredVars))
	for k := range c.registeredVars {
		ret = append(ret, k)
	}
	return ret
}

// Var returns an expvar variable that shows all the current configuration variables and their
// current value
func (c *Distconf) Var() expvar.Var {
	return expvar.Func(func() interface{} {
		c.varsMutex.Lock()
		defer c.varsMutex.Unlock()

		m := make(map[string]interface{})
		for name, v := range c.registeredVars {
			m[name] = v.distvar.GenericGet()
		}
		return m
	})
}

// Int object that can be referenced to get integer values from a backing config
func (c *Distconf) Int(key string, defaultVal int64) *Int {
	s := &intConf{
		defaultVal: defaultVal,
		Int: Int{
			currentVal: defaultVal,
		},
	}
	// Note: in race conditions 's' may not be the thing actually returned
	ret, okCast := c.createOrGet(key, s).(*intConf)
	if !okCast {
		c.logger(key, nil, "Registering key with multiple types!  FIX ME!!!!")
		return nil
	}
	return &ret.Int
}

// Float object that can be referenced to get float values from a backing config
func (c *Distconf) Float(key string, defaultVal float64) *Float {
	s := &floatConf{
		defaultVal: defaultVal,
		Float: Float{
			currentVal: math.Float64bits(defaultVal),
		},
	}
	// Note: in race conditions 's' may not be the thing actually returned
	ret, okCast := c.createOrGet(key, s).(*floatConf)
	if !okCast {
		c.logger(key, nil, "Registering key with multiple types!  FIX ME!!!!")
		return nil
	}
	return &ret.Float
}

// Struct object that can be referenced to decode a json representation of an instance
// of the type that is provided in the defaultVal. The DefaultValue provided sets the type that
// the backend josn will be unmarshelled into.
func (c *Distconf) Struct(key string, defaultVal interface{}) *Struct {
	s := &structConf{
		defaultVal: defaultVal,
		t:          reflect.TypeOf(defaultVal),
	}
	s.currentVal.Store(defaultVal)

	ret, okCast := c.createOrGet(key, s).(*structConf)
	// check that the new default value has the same type as the previous default value.
	if !okCast || ret.t != s.t {
		c.logger(key, nil, "Registering key with multiple types!  FIX ME!!!!")
		return nil
	}
	return &ret.Struct
}

// Str object that can be referenced to get string values from a backing config
func (c *Distconf) Str(key string, defaultVal string) *Str {
	s := &strConf{
		defaultVal: defaultVal,
	}
	s.currentVal.Store(defaultVal)
	// Note: in race conditions 's' may not be the thing actually returned
	ret, okCast := c.createOrGet(key, s).(*strConf)
	if !okCast {
		c.logger(key, nil, "Registering key with multiple types!  FIX ME!!!!")
		return nil
	}
	return &ret.Str
}

// Bool object that can be referenced to get boolean values from a backing config
func (c *Distconf) Bool(key string, defaultVal bool) *Bool {
	var defautlAsInt int32
	if defaultVal {
		defautlAsInt = 1
	} else {
		defautlAsInt = 0
	}

	s := &boolConf{
		defaultVal: defautlAsInt,
		Bool: Bool{
			currentVal: defautlAsInt,
		},
	}
	// Note: in race conditions 's' may not be the thing actually returned
	ret, okCast := c.createOrGet(key, s).(*boolConf)
	if !okCast {
		c.logger(key, nil, "Registering key with multiple types!  FIX ME!!!!")
		return nil
	}
	return &ret.Bool
}

// Duration returns a duration object that calls ParseDuration() on the given key
func (c *Distconf) Duration(key string, defaultVal time.Duration) *Duration {
	s := &durationConf{
		defaultVal: defaultVal,
		Duration: Duration{
			currentVal: defaultVal.Nanoseconds(),
		},
		logger: c.logger,
	}
	// Note: in race conditions 's' may not be the thing actually returned
	ret, okCast := c.createOrGet(key, s).(*durationConf)
	if !okCast {
		c.logger(key, nil, "Registering key with multiple types!  FIX ME!!!!")
		return nil
	}
	return &ret.Duration
}

// Close this config framework's readers.  Config variable results are undefined after this call.
func (c *Distconf) Close() {
	c.varsMutex.Lock()
	defer c.varsMutex.Unlock()
	for _, backing := range c.Readers {
		backing.Close()
	}
}

func (c *Distconf) refresh(key string, configVar configVariable) bool {
	dynamicReadersOnPath := false
	for _, backing := range c.Readers {
		if !dynamicReadersOnPath {
			_, ok := backing.(Dynamic)
			if ok {
				dynamicReadersOnPath = true
			}
		}

		v, e := backing.Get(key)
		if e != nil {
			c.logger(key, e, "Unable to read from backing")
			continue
		}
		if v != nil {
			e = configVar.Update(v)
			if e != nil {
				c.logger(key, e, "Invalid config bytes")
			}
			return dynamicReadersOnPath
		}
	}

	e := configVar.Update(nil)
	if e != nil {
		c.logger(key, e, "Unable to set bytes to nil/clear")
	}

	// If this is false, then the variable is fixed and can never change
	return dynamicReadersOnPath
}

func (c *Distconf) watch(key string, configVar configVariable) {
	for _, backing := range c.Readers {
		d, ok := backing.(Dynamic)
		if ok {
			err := d.Watch(key, c.onBackingChange)
			if err != nil {
				c.logger(key, err, "Unable to watch for config var")
			}
		}
	}
}

func (c *Distconf) createOrGet(key string, defaultVar configVariable) configVariable {
	c.varsMutex.Lock()
	rv, exists := c.registeredVars[key]
	if !exists {
		rv = &registeredVariableTracker{
			distvar: defaultVar,
		}
		if c.registeredVars == nil {
			c.registeredVars = make(map[string]*registeredVariableTracker)
		}
		c.registeredVars[key] = rv
	}
	c.varsMutex.Unlock()

	rv.hasInitialized.Do(func() {
		dynamicOnPath := c.refresh(key, rv.distvar)
		if dynamicOnPath {
			c.watch(key, rv.distvar)
		}
	})
	return rv.distvar
}

func (c *Distconf) onBackingChange(key string) {
	c.varsMutex.Lock()
	m, exists := c.registeredVars[key]
	c.varsMutex.Unlock()
	if !exists {
		c.logger(key, nil, "Backing callback on variable that doesn't exist")
		return
	}
	c.refresh(key, m.distvar)
}

// Reader can get a []byte value for a config key
type Reader interface {
	// Get should return the given config value.  If the value does not exist, it should return nil, nil.
	Get(key string) ([]byte, error)
	Close()
}

// Writer can modify Config properties
type Writer interface {
	Write(key string, value []byte) error
}

// A Dynamic config can change what it thinks a value is over time.
type Dynamic interface {
	// Watch should execute callback function whenever the key changes.  The parameter to callback should be the
	// key's name.
	Watch(key string, callback func(string)) error
}

// A ReaderWriter can both read and write configuration information
type ReaderWriter interface {
	Reader
	Writer
}
