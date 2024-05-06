package log_test

import (
	"testing"

	"github.com/chelnak/gh-changelog/internal/log"
	"github.com/stretchr/testify/require"
)

func TestGetLoggerTypeReturnsTheCorrectType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     log.LoggerType
		hasError bool
	}{
		{
			name:     "returns a console logger",
			input:    "console",
			want:     log.CONSOLE,
			hasError: false,
		},
		{
			name:     "returns a spinner logger",
			input:    "spinner",
			want:     log.SPINNER,
			hasError: false,
		},
		{
			name:     "returns an error for an invalid logger",
			input:    "invalid",
			want:     log.SPINNER,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerType := log.GetLoggerType(tt.input)
			require.Equal(t, tt.want, loggerType)
		})
	}
}
