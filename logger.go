package suppress

import (
	"errors"

	tlog "github.com/Dendromeda/suppresserrors/log"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
)

func AddSuppressedErrorTypes(logger *log.Logger, suppressedErrorTypes ...string) log.Logger {
	return &Logger{
		internalLogger:       *logger,
		suppressedErrorTypes: suppressedErrorTypes,
	}
}

func NewLoggerWithSuppressedTypes(suppressedErrorTypes ...string) log.Logger {
	return &Logger{
		internalLogger:       tlog.NewDefaultLogger(),
		suppressedErrorTypes: suppressedErrorTypes,
	}
}

type Logger struct {
	internalLogger       log.Logger
	suppressedErrorTypes []string
}

// Debug implements log.Logger.
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.internalLogger.Debug(msg, keyvals...)
}

// Error implements log.Logger.
func (l *Logger) Error(msg string, keyvals ...interface{}) {
	//x, _ := json.MarshalIndent(keyvals, " ", "    ")

	if len(keyvals) >= 16 {

		switch t := keyvals[15].(type) {
		case error:
			var appErr *temporal.ApplicationError
			if errors.As(t, &appErr) {
				for _, e := range l.suppressedErrorTypes {
					if appErr.Type() == e {
						l.internalLogger.Info(msg, keyvals...)
						return
					}
				}
			}
		}
	}
	l.internalLogger.Error(msg, keyvals...)
}

// Info implements log.Logger.
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.internalLogger.Info(msg, keyvals...)
}

// Warn implements log.Logger.
func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.internalLogger.Warn(msg, keyvals...)
}

var _ log.Logger = &Logger{}
