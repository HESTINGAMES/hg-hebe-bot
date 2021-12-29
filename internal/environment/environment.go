package environment

import (
	"os"
)

// Environment constants
const (
	prod     = "prod"
	dev      = "dev"
	localDev = "localDev"
)

var env string

func init() {
	env = Environment()
}

func Environment() string {
	if env != "" {
		return env
	}
	env = os.Getenv("ENVIRONMENT")
	if env == "" {
		env = localDev
	}
	return env
}

func IsDev() bool {
	return Environment() == dev
}

func IsProd() bool {
	return Environment() == prod
}

func IsLocal() bool {
	return Environment() == localDev
}
