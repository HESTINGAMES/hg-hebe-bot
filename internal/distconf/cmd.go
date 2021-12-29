package distconf

import (
	"os"
	"strings"
)

// CommandLine gets distconf values from the command line
type CommandLine struct {
	Prefix string
	Source []string
}

func (p *CommandLine) source() []string {
	if p.Source == nil {
		return os.Args
	}
	return p.Source
}

// Get looks for "prefix+key=xyz"
func (p *CommandLine) Get(key string) ([]byte, error) {
	argPrefix := p.Prefix + key + "="
	for _, arg := range p.source() {
		if !strings.HasPrefix(arg, argPrefix) {
			continue
		}
		argSuffix := arg[len(argPrefix):]
		return []byte(argSuffix), nil
	}
	return nil, nil
}

// Close does nothing
func (p *CommandLine) Close() {
}
