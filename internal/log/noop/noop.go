// Package noop is responsible for providing a no-op logger
package noop

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (l *NoopLogger) Infof(string, ...interface{}) {}

func (l *NoopLogger) Errorf(string, ...interface{}) {}

func (l *NoopLogger) Complete() {}
