package logs

import (
	"fmt"
	golog "log"

	"github.com/hestingames/hg-hebe-bot/internal/environment"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	initialSampleSize = 100
	sampleRate        = 100
)

// Logger wraps the implementation of the logger to provide a more flexible
// interface.
type Logger struct {
	*zap.Logger
}

// New constructs a new instance of a Logger.
func New() (*Logger, error) {
	var log *zap.Logger
	var err error

	if environment.IsLocal() {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		log, err = config.Build()
	} else {
		// Assume deployed in some environment like staging, canary, or production.
		config := zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Encoding:          "console", // console, json
			DisableCaller:     true,
			DisableStacktrace: true,
			EncoderConfig:     zap.NewProductionEncoderConfig(),
			OutputPaths:       []string{"stderr"},
			ErrorOutputPaths:  []string{"stderr"},

			// Zap samples by logging the first N entries with a given level and
			// message each tick. If more Entries with the same level and message
			// are seen during the same interval, every Mth message is logged and the
			// rest are dropped.
			Sampling: &zap.SamplingConfig{
				Initial:    initialSampleSize, // accept the first `Initial` logs each second...
				Thereafter: sampleRate,        // and one log in every `Thearafter` after that...
			},
		}
		log, err = config.Build()
	}

	if err != nil {
		return nil, err
	}

	// Set log as global
	zap.ReplaceGlobals(log)
	defer log.Sync() // flushes buffer, if any

	return &Logger{Logger: log}, nil
}

func (l *Logger) With(fields ...zapcore.Field) *Logger {
	return &Logger{l.Logger.With(fields...)}
}

func (l *Logger) NewStdLogger() *golog.Logger {
	return zap.NewStdLog(l.Logger)
}

// Log allows zap logger to meet the Fulton TwitchLogging interface: code.justin.tv/amzn/TwitchLogging
// Log does a best effort logging of the given args using zap as an underlying logger
// The fulton logger expects a main message, then a series of keys and values as
// followup arguments for structured logging, like
//     fultonlogger.Log("my main message!", "KEY1", "VALUE1", "KEY2", "VALUE2")
//
// Here, we will usher this into something that zap likes by taking those keypairs
// and shoving them into zap.String() calls
//
// Implementation heavily adapted from JSON fulton logger:
// https://git-aws.internal.justin.tv/amzn/TwitchLoggingCommonLoggers/blob/mainline/json_logger.go#L35-L54
func (l *Logger) Log(msg string, keyvals ...interface{}) {
	numFields := ((len(keyvals) + (len(keyvals) % 2)) / 2)
	fields := make([]zapcore.Field, numFields)
	for i := 0; i < len(keyvals); i += 2 {
		var key, value interface{}
		if i == len(keyvals)-1 {
			key, value = "UNKEYED_VALUE", keyvals[i]
		} else {
			key, value = keyvals[i], keyvals[i+1]
		}
		fields[i/2] = zap.String(fmt.Sprint(key), fmt.Sprint(value))
	}
	l.Info(msg, fields...)
}
