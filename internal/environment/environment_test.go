package environment_test

import (
	"os"
	"testing"

	"github.com/chelnak/gh-changelog/internal/environment"
	"github.com/stretchr/testify/assert"
)

func TestIsCIReturnsTrueWhenRunningInCI(t *testing.T) {
	_ = os.Setenv("CI", "true")
	defer func() {
		_ = os.Unsetenv("CI")
	}()
	assert.True(t, environment.IsCI())
}
