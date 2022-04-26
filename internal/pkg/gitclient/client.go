package gitclient

import (
	"os/exec"
	"strings"
	"time"
)

type Tag struct {
	Name string
	Sha  string
	Date time.Time
}

type execOptions struct {
	args []string
}

type Git struct {
}

func (g *Git) exec(opts execOptions) (string, error) {
	// TODO: Consider not using a private exec function and hardcode
	// each call to git in the respective command.
	// For now, the lint check is disabled.
	output, err := exec.Command("git", opts.args...).Output() // #nosec G204
	if err != nil {
		return "", err
	}

	return strings.Trim(string(output), "\n"), nil
}

func (g *Git) GetFirstCommit() (string, error) {
	return g.exec(execOptions{
		args: []string{"rev-list", "--max-parents=0", "HEAD"},
	})
}

func (g *Git) GetDateOfHash(hash string) (time.Time, error) {
	date, err := g.exec(execOptions{
		args: []string{"log", "-1", "--format=%cI", hash, "--date=local"},
	})

	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(time.RFC3339, date, time.Local)
}

func NewGitHandler() *Git {
	return &Git{}
}
