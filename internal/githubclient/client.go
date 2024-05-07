// Package githubclient is a wrapper around the githubv4 client.
// It's purpose is to provide abstraction for some graphql queries
// that retrieve data for the changelog.
package githubclient

import (
	"context"
	"fmt"
	"time"

	"github.com/chelnak/gh-changelog/internal/utils"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/shurcooL/githubv4"
)

type repoContext struct {
	owner string
	name  string
}

type GitHubClient interface {
	GetTags() ([]Tag, error)
	GetPullRequestsBetweenDates(from, to time.Time) ([]PullRequest, error)
	GetRepoName() string
	GetRepoOwner() string
}

type githubClient struct {
	base        *githubv4.Client
	repoContext repoContext
	httpContext context.Context
}

func (client *githubClient) GetRepoName() string {
	return client.repoContext.name
}

func (client *githubClient) GetRepoOwner() string {
	return client.repoContext.owner
}

func NewGitHubClient() (GitHubClient, error) {
	httpClient, err := api.DefaultHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("could not create initial client: %s", err)
	}

	g := githubv4.NewClient(httpClient)

	currentRepository, err := utils.GetRepoContext()
	if err != nil {
		return nil, err
	}

	client := &githubClient{
		base: g,
		repoContext: repoContext{
			owner: currentRepository.Owner,
			name:  currentRepository.Name,
		},
		httpContext: context.Background(),
	}

	return client, nil
}
