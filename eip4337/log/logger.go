package log

import "github.com/tendermint/tendermint/libs/log"

type Logger = log.Logger

type NoOpLogger struct{}

func (logger *NoOpLogger) Debug(msg string, keyvals ...interface{}) {}

func (logger *NoOpLogger) Info(msg string, keyvals ...interface{}) {}

func (logger *NoOpLogger) Error(msg string, keyvals ...interface{}) {}

func (logger *NoOpLogger) With(keyvals ...interface{}) Logger { return logger }

var _ Logger = (*NoOpLogger)(nil)

func EnsureLogger(logger Logger) Logger {
	if logger == nil {
		return &NoOpLogger{}
	}
	return logger
}
