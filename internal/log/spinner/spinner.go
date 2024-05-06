// Package spinner is responsible for providing a spinner logger
package spinner

import (
	"fmt"

	"github.com/chelnak/ysmrr"
)

type SpinnerLogger struct {
	manager ysmrr.SpinnerManager
	spinner *ysmrr.Spinner
}

func NewSpinnerLogger() *SpinnerLogger {
	manager := ysmrr.NewSpinnerManager()
	spinner := manager.AddSpinner("Loading..")
	manager.Start()
	return &SpinnerLogger{
		manager: manager,
		spinner: spinner,
	}
}

func (l *SpinnerLogger) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.spinner.UpdateMessage(message)
}

func (l *SpinnerLogger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.spinner.UpdateMessage(message)
	l.spinner.Error()
}

func (l *SpinnerLogger) Complete() {
	l.spinner.Complete()
	l.manager.Stop()
}
