package githubclient

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/repository"
)

// RepoContext is a struct that contains the current repository owner and name.
type RepoContext struct {
	Owner string
	Name  string
}

// GetRepoContext returns a new RepoContext struct with the current repository owner and name.
func GetRepoContext() (RepoContext, error) {
	currentRepository, err := repository.Current()
	if err != nil {
		if strings.Contains(err.Error(), "not a git repository (or any of the parent directories)") {
			return RepoContext{}, fmt.Errorf("the current directory is not a git repository")
		}

		return RepoContext{}, err
	}

	return RepoContext{
		Owner: currentRepository.Owner,
		Name:  currentRepository.Name,
	}, nil
}
