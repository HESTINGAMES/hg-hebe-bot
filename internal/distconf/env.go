package distconf

import "os"

// Env allows fetching configuration from the env
type Env struct {
	// Prefix is a forced prefix in front of config variables
	Prefix string
	// OsGetenv is a stub for os.Getenv()
	OsGetenv func(string) string
}

func (p *Env) osGetenv(key string) string {
	if p.OsGetenv != nil {
		return p.OsGetenv(key)
	}
	return os.Getenv(key)
}

// Get prefix + key from env if it exists.
func (p *Env) Get(key string) ([]byte, error) {
	val := p.osGetenv(p.Prefix + key)
	if val == "" {
		return nil, nil
	}
	return []byte(val), nil
}

// Close does nothing and exists just to satisfy an interface
func (p *Env) Close() {
}
