// Package logging provides a simple logging interface for the
// application.
package logging

import (
	"fmt"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/environment"
)

// LoggerType is an enum for the different types of logger
type LoggerType int64

const (
	TEXT LoggerType = iota
	SPINNER
)

// Logger is the interface for logging in the application.
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Complete()
	GetType() LoggerType
}

// NewLogger returns a new logger based on the type passed in.
func NewLogger(loggerType LoggerType) Logger {
	switch loggerType {
	case TEXT:
		return newTextLogger()
	case SPINNER:
		return newSpinnerLogger()
	default:
		return newSpinnerLogger()
	}
}

// GetLoggerType returns the logger type from the string value
// passed in. This is a convenience function for the CLI.
func GetLoggerType(name string) (LoggerType, error) {
	if name == "" {
		name = configuration.Config.Logger

		// If we're running in a CI environment then we don't want to
		// use the spinner.
		if environment.IsCI() {
			name = "text"
		}
	}

	var loggerType LoggerType
	switch name {
	case "text":
		loggerType = TEXT
	case "spinner":
		loggerType = SPINNER
	default:
		return loggerType, fmt.Errorf("'%s' is not a valid logger. Valid values are 'spinner' and 'text'", name)
	}

	return loggerType, nil
}
