// Package githubclient is a wrapper around api.DefaultGraphQLClient.
// Its purpose is to provide abstraction for some graphql queries
// that retrieve data for the changelog.
package githubclient

import (
	"context"
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
)

type repoContext struct {
	owner string
	name  string
}

type GitHub struct {
	base        *api.GraphQLClient
	repoContext repoContext
	httpContext context.Context
}

func (client *GitHub) GetRepoName() string {
	return client.repoContext.name
}

func (client *GitHub) GetRepoOwner() string {
	return client.repoContext.owner
}

func NewGitHubClient() (*GitHub, error) {
	g, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("could not create graphql client: %w", err)
	}

	currentRepository, err := GetRepoContext()
	if err != nil {
		return nil, err
	}

	client := &GitHub{
		base: g,
		repoContext: repoContext{
			owner: currentRepository.Owner,
			name:  currentRepository.Name,
		},
		httpContext: context.Background(),
	}

	return client, nil
}
