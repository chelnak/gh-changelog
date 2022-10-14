package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type textLogger struct {
	loggerType LoggerType
}

func newTextLogger() Logger {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &textLogger{
		loggerType: TEXT,
	}
}

func (t *textLogger) Infof(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

func (t *textLogger) Errorf(format string, args ...interface{}) {
	log.Error().Msgf(format, args...)
}

func (t *textLogger) Complete() {}

func (t *textLogger) GetType() LoggerType {
	return t.loggerType
}
