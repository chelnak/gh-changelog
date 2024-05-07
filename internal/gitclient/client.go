// Package gitclient is responsible for providing an interface
// to the local git binary. It provides predefined calls that can
// be easily consumed by other packages.
package gitclient

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type execContext = func(name string, arg ...string) *exec.Cmd

type execOptions struct {
	args []string
}

type Git struct {
	execContext execContext
}

func (g Git) exec(opts execOptions) (string, error) {
	// TODO: Consider not using a private exec function and hardcode
	// each call to git in the respective command.
	// For now, the lint check is disabled.
	// output, err := exec.Command("git", opts.args...).Output() // #nosec G204
	output, err := g.execContext("git", opts.args...).Output() // #nosec G204
	if err != nil {
		return "", fmt.Errorf("git command failed: %s\n%s", strings.Join(opts.args, " "), err)
	}

	return strings.Trim(string(output), "\n"), nil
}

func (g Git) GetFirstCommit() (string, error) {
	response, err := g.exec(execOptions{
		args: []string{"rev-list", "--max-parents=0", "HEAD", "--reverse"},
	})

	if err != nil {
		return "", err
	}

	hashes := strings.Split(response, "\n")

	// if len(hashes) > 1 {
	// 	//If we arrive here it means that rev-list has returned more than one commit.
	// 	//This can happen when there are orphaned commits in the repository.
	// 	//We split the response by newline and return the the item at position 0.
	// 	//TODO: Logging should be added here to explain the situation.
	// }

	return hashes[0], nil
}

func (g Git) GetLastCommit() (string, error) {
	return g.exec(execOptions{
		args: []string{"log", "-1", "--pretty=format:%H"},
	})
}

func (g Git) GetDateOfHash(hash string) (time.Time, error) {
	date, err := g.exec(execOptions{
		args: []string{"log", "-1", "--format=%cI", hash, "--date=local"},
	})

	if err != nil {
		return time.Time{}, err
	}

	return time.ParseInLocation(time.RFC3339, date, time.Local)
}

func (g Git) GetCurrentBranch() (string, error) {
	return g.exec(execOptions{
		args: []string{"rev-parse", "--abbrev-ref", "HEAD"},
	})
}

func (g Git) FetchAll() error {
	_, err := g.exec(execOptions{
		args: []string{"fetch", "--all"},
	})

	return err
}

func (g Git) Tags() ([]string, error) {
	tags, err := g.exec(execOptions{
		args: []string{"for-each-ref", "--format='%(refname:short) %(objectname) %(taggerdate:iso-strict)'", "refs/tags"},
	})

	if err != nil {
		return nil, err
	}
	t := strings.Split(tags, "\n")
	return t, nil
}

func (g Git) IsAncestorOf(commit, ancestor string) (bool, error) {
	_, err := g.exec(execOptions{
		args: []string{"merge-base", "--is-ancestor", commit, ancestor},
	})

	if err != nil {
		if strings.Contains(err.Error(), "exit status 1") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func NewGitClient(cmdContext execContext) Git {
	return Git{
		execContext: cmdContext,
	}
}
