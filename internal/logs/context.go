package logs

import (
	"context"
	golog "log"
	"os"
)

type key string

var contextKey = key("logger")

func (l Logger) AddToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, l)
}

func FromContext(ctx context.Context) *Logger {
	if v := ctx.Value(contextKey); v != nil {
		if log, ok := v.(Logger); ok {
			return &log
		}
	}

	log, err := New()
	if err != nil {
		golog.Fatal("could not create new logger for environment", os.Environ())
	}

	return log
}
