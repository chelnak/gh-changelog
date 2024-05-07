package update_test

import (
	"testing"

	"github.com/chelnak/gh-changelog/internal/update"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestCheckForUpdates(t *testing.T) {
	tests := []struct {
		name           string
		currentVersion string
		nextVersion    string
		want           bool
	}{
		{
			name:           "an update is available",
			currentVersion: "changelog version 1.0.0",
			nextVersion:    "v1.0.1",
			want:           true,
		},
		{
			name:           "no update is available",
			currentVersion: "changelog version 1.0.0",
			nextVersion:    "1.0.0",
			want:           false,
		},
	}

	defer gock.Off()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gock.New("https://api.github.com").
				Get("/repos/chelnak/gh-changelog/releases/latest").
				Reply(200).
				JSON(map[string]string{"tag_name": tt.nextVersion})

			got := update.CheckForUpdate(tt.currentVersion)
			require.Equal(t, tt.want, got)
		})
	}
}
