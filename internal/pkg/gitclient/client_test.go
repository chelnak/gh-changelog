package gitclient_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/chelnak/gh-changelog/internal/pkg/gitclient"
	"github.com/stretchr/testify/assert"
)

var testStdoutValue = "test"

func fakeExecSuccess(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestShellProcessSuccess", "--", command}
	cs = append(cs, args...)
	expectedOutput := os.Getenv("GO_TEST_PROCESS_EXPECTED_OUTPUT")
	cmd := exec.Command(os.Args[0], cs...) // #nosec
	cmd.Env = []string{"GO_TEST_PROCESS=1", fmt.Sprintf("GO_TEST_PROCESS_EXPECTED_OUTPUT=%s", expectedOutput)}
	return cmd
}

func fakeExecFailure(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestShellProcessFailure", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...) // #nosec
	cmd.Env = []string{"GO_TEST_PROCESS=1"}
	return cmd
}

func safeSetMockOutput(output string) func() {
	_ = os.Setenv("GO_TEST_PROCESS_EXPECTED_OUTPUT", output)
	return func() {
		_ = os.Unsetenv("GO_TEST_PROCESS_EXPECTED_OUTPUT")
	}
}

func TestShellProcessSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	output := testStdoutValue
	if os.Getenv("GO_TEST_PROCESS_EXPECTED_OUTPUT") != "" {
		output = os.Getenv("GO_TEST_PROCESS_EXPECTED_OUTPUT")
	}

	fmt.Fprint(os.Stdout, output)
	os.Exit(0)
}

func TestShellProcessFailure(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}
	fmt.Fprint(os.Stderr, "error")
	os.Exit(1)
}

func TestGetFirstCommitSuccess(t *testing.T) {
	gitClient := gitclient.NewGitClient(fakeExecSuccess)
	commit, err := gitClient.GetFirstCommit()

	assert.NoError(t, err)
	assert.Equal(t, "test", commit)
}

func TestGetFirstCommitFailure(t *testing.T) {
	gitClient := gitclient.NewGitClient(fakeExecFailure)
	_, err := gitClient.GetFirstCommit()

	assert.Error(t, err)
}

func TestGetFirstCommitWithOrphansSuccess(t *testing.T) {
	defer safeSetMockOutput("test-hash-0\ntest-hash-1")()

	gitClient := gitclient.NewGitClient(fakeExecSuccess)
	commit, err := gitClient.GetFirstCommit()

	assert.NoError(t, err)
	assert.Equal(t, "test-hash-0", commit)
}

func TestGetLastCommitSuccess(t *testing.T) {
	gitClient := gitclient.NewGitClient(fakeExecSuccess)
	commit, err := gitClient.GetLastCommit()

	assert.NoError(t, err)
	assert.Equal(t, "test", commit)
}

func TestGetLastCommitFailure(t *testing.T) {
	gitClient := gitclient.NewGitClient(fakeExecFailure)
	_, err := gitClient.GetLastCommit()

	assert.Error(t, err)
}

func TestGetDateOfHashSuccess(t *testing.T) {
	mockDate := "2022-04-18T19:31:31+00:00"
	defer safeSetMockOutput(mockDate)()

	gitClient := gitclient.NewGitClient(fakeExecSuccess)
	date, err := gitClient.GetDateOfHash("test-hash")
	expectedDate, _ := time.ParseInLocation(time.RFC3339, mockDate, time.Local)

	assert.NoError(t, err)
	assert.Equal(t, expectedDate, date)
}
