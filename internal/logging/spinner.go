package logging

import (
	"fmt"

	"github.com/chelnak/ysmrr"
)

type spinnerLogger struct {
	manager    ysmrr.SpinnerManager
	spinner    *ysmrr.Spinner
	loggerType LoggerType
}

func newSpinnerLogger() Logger {
	manager := ysmrr.NewSpinnerManager()
	spinner := manager.AddSpinner("Loading..")
	manager.Start()
	return &spinnerLogger{
		manager:    manager,
		spinner:    spinner,
		loggerType: SPINNER,
	}
}

func (s *spinnerLogger) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	s.spinner.UpdateMessage(message)
}

func (s *spinnerLogger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	s.spinner.UpdateMessage(message)
	s.spinner.Error()
}

func (s *spinnerLogger) Complete() {
	s.spinner.Complete()
	s.manager.Stop()
}

func (s *spinnerLogger) GetType() LoggerType {
	return s.loggerType
}
