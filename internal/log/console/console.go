// Package console is responsible for providing a console logger
package console

import (
	"os"

	"github.com/rs/zerolog"
	console "github.com/rs/zerolog/log"
)

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	console.Logger = console.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Infof(format string, args ...interface{}) {
	console.Info().Msgf(format, args...)
}

func (l *ConsoleLogger) Errorf(format string, args ...interface{}) {
	console.Error().Msgf(format, args...)
}

func (l *ConsoleLogger) Complete() {}
