package log

type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})

	With(keyvals ...interface{}) Logger
}

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
