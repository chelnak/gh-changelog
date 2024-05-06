// Package log provides a simple logging interface for the
// application.
package log

import (
	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/environment"
	"github.com/chelnak/gh-changelog/internal/log/console"
	"github.com/chelnak/gh-changelog/internal/log/noop"
	"github.com/chelnak/gh-changelog/internal/log/spinner"
)

// LoggerType is an enum for the different types of logger
type LoggerType int64

// Logger is the interface for logging in the application.
type logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Complete()
}

const (
	CONSOLE LoggerType = iota
	SPINNER
	NOOP
)

var base logger

// SetupLogging returns a new logger based on the type passed in.
func SetupLogging(loggerType LoggerType) {
	switch loggerType {
	case CONSOLE:
		base = console.NewConsoleLogger()
	case SPINNER:
		base = spinner.NewSpinnerLogger()
	case NOOP:
		base = noop.NewNoopLogger()
	default:
		base = spinner.NewSpinnerLogger()
	}
}

// Infof logs an info message.
func Infof(format string, args ...interface{}) {
	base.Infof(format, args...)
}

// Errorf logs an error message.
func Errorf(format string, args ...interface{}) {
	base.Errorf(format, args...)
}

// Complete completes the logging process.
func Complete() {
	base.Complete()
}

// GetLoggerType returns the logger type from the string value
// passed in. This is a convenience function for the CLI.
func GetLoggerType(name string) LoggerType {
	if name == "" {
		name = configuration.Config.Logger

		// If we're running in a CI environment then we don't want to
		// use the spinner.
		if environment.IsCI() {
			name = "console"
		}
	}

	var loggerType LoggerType
	switch name {
	case "console":
		loggerType = CONSOLE
	case "spinner":
		loggerType = SPINNER
	case "noop":
		loggerType = NOOP
	default:
		loggerType = SPINNER
	}

	return loggerType
}
