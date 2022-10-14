package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type consoleLogger struct {
	loggerType LoggerType
}

func newConsoleLogger() Logger {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &consoleLogger{
		loggerType: CONSOLE,
	}
}

func (c *consoleLogger) Infof(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

func (c *consoleLogger) Errorf(format string, args ...interface{}) {
	log.Error().Msgf(format, args...)
}

func (c *consoleLogger) Complete() {}

func (c *consoleLogger) GetType() LoggerType {
	return c.loggerType
}
