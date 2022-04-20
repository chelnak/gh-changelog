package githubclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cli/go-gh"
	"github.com/google/go-github/v43/github"
)

type RepoContext struct {
	Owner string
	Name  string
}

type GitHubClient struct {
	baseClient  *github.Client
	RepoContext RepoContext
	httpContext context.Context
}

// NewGitHubClient returns an instance of *github.Client that has
// been initialised with the HttpClient provided by the go-gh module.
// This ensures that the client is compatible with the GitHub CLI.
func NewGitHubClient() (*GitHubClient, error) {
	httpClient, err := gh.HTTPClient(nil)
	if err != nil {
		return nil, err
	}

	g := github.NewClient(httpClient)

	currentRepository, err := gh.CurrentRepository()
	if err != nil {
		if strings.Contains(err.Error(), "not a git repository (or any of the parent directories)") {
			return nil, fmt.Errorf("the current directory is not a git repository")
		}

		return nil, err
	}

	client := &GitHubClient{
		baseClient: g,
		RepoContext: RepoContext{
			Owner: currentRepository.Owner(),
			Name:  currentRepository.Name(),
		},
		httpContext: context.Background(),
	}

	return client, nil
}

func (c *GitHubClient) GetPullRequestsBetweenDates(fromDate time.Time, toDate time.Time) ([]*github.Issue, error) {
	searchOptions := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	query := fmt.Sprintf(
		"repo:%s/%s type:pr is:merged merged:%s..%s",
		c.RepoContext.Owner, c.RepoContext.Name,
		fromDate.Format("2006-01-02T15:04:05+00:00"), toDate.Format("2006-01-02T15:04:05+00:00"),
	)
	result, response, err := c.baseClient.Search.Issues(c.httpContext, query, &searchOptions)

	// Feels gross to be doing recursion... but maybe it might be the right idea??
	// could memoize retry and set some sensible limit??
	if response.Rate.Remaining == 0 {
		now := time.Now()
		resetTime := response.Rate.Reset.Time.Add(time.Second * time.Duration(5))
		timeToWait := resetTime.Sub(now)

		time.Sleep(timeToWait)

		return c.GetPullRequestsBetweenDates(fromDate, toDate)
	}

	if err != nil {
		return nil, err
	}

	return result.Issues, nil
}
