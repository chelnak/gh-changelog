package environment_test

import (
	"os"
	"testing"

	"github.com/chelnak/gh-changelog/internal/environment"
	"github.com/stretchr/testify/require"
)

func TestIsCIReturnsTrueWhenRunningInCI(t *testing.T) {
	_ = os.Setenv("CI", "true")
	defer func() {
		_ = os.Unsetenv("CI")
	}()
	require.True(t, environment.IsCI())
}
