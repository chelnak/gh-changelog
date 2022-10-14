package logging_test

import (
	"testing"

	"github.com/chelnak/gh-changelog/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestNewLoggerReturnsTheCorrectType(t *testing.T) {
	tests := []struct {
		name     string
		want     logging.LoggerType
		hasError bool
	}{
		{
			name:     "returns a text logger",
			want:     logging.TEXT,
			hasError: false,
		},
		{
			name:     "returns a spinner logger",
			want:     logging.SPINNER,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := logging.NewLogger(tt.want)
			assert.Equal(t, tt.want, got.GetType())
		})
	}
}

func TestGetLoggerTypeReturnsTheCorrectType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     logging.LoggerType
		hasError bool
	}{
		{
			name:     "returns a text logger",
			input:    "text",
			want:     logging.TEXT,
			hasError: false,
		},
		{
			name:     "returns a spinner logger",
			input:    "spinner",
			want:     logging.SPINNER,
			hasError: false,
		},
		{
			name:     "returns an error for an invalid logger",
			input:    "invalid",
			want:     logging.SPINNER,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := logging.GetLoggerType(tt.input)
			if tt.hasError {
				assert.Error(t, err)
				assert.Equal(t, "'invalid' is not a valid logger. Valid values are 'spinner' and 'text'", err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
