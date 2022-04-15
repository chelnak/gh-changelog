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
	Repo  string
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
			Repo:  currentRepository.Name(),
		},
		httpContext: context.Background(),
	}

	return client, nil
}

func (c *GitHubClient) GetPullRequestsBeforeDate(date time.Time) ([]*github.Issue, error) {
	searchOptions := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	query := fmt.Sprintf(
		"repo:%s/%s type:pr is:merged merged:<%s",
		c.RepoContext.Owner, c.RepoContext.Repo,
		date.Format("2006-01-02T15:04:05+00:00"),
	)

	result, _, err := c.baseClient.Search.Issues(c.httpContext, query, &searchOptions)
	if err != nil {
		return nil, err
	}

	return result.Issues, nil
}

func (c *GitHubClient) GetPullRequestAfterDate(date time.Time) ([]*github.Issue, error) {
	searchOptions := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	query := fmt.Sprintf(
		"repo:%s/%s type:pr is:merged merged:>%s",
		c.RepoContext.Owner, c.RepoContext.Repo,
		date.Format("2006-01-02T15:04:05+00:00"),
	)
	result, _, err := c.baseClient.Search.Issues(c.httpContext, query, &searchOptions)
	if err != nil {
		return nil, err
	}

	return result.Issues, nil
}

func (c *GitHubClient) GetPullRequestsBetweenDates(fromDate time.Time, toDate time.Time) ([]*github.Issue, error) {
	searchOptions := github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	query := fmt.Sprintf(
		"repo:%s/%s type:pr is:merged merged:%s..%s",
		c.RepoContext.Owner, c.RepoContext.Repo,
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

func (c *GitHubClient) GetTag(name string) (*github.RepositoryTag, error) {
	listTagOptions := github.ListOptions{
		PerPage: 100,
	}
	tags, _, err := c.baseClient.Repositories.ListTags(
		c.httpContext,
		c.RepoContext.Owner,
		c.RepoContext.Repo,
		&listTagOptions,
	)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if strings.EqualFold(tag.GetName(), name) {
			return tag, nil
		}
	}

	return nil, fmt.Errorf("tag %s not found", name)
}

func (c *GitHubClient) GetTags() ([]*github.RepositoryTag, error) {
	listTagsOptions := github.ListOptions{
		PerPage: 100,
	}
	tags, _, err := c.baseClient.Repositories.ListTags(
		c.httpContext,
		c.RepoContext.Owner,
		c.RepoContext.Repo,
		&listTagsOptions,
	)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (c *GitHubClient) GetCommit(sha string) (*github.Commit, error) {
	commit, _, err := c.baseClient.Git.GetCommit(
		c.httpContext,
		c.RepoContext.Owner,
		c.RepoContext.Repo,
		sha,
	)
	if err != nil {
		return nil, err
	}
	return commit, nil
}

func (c *GitHubClient) GetFirstTag() (*github.RepositoryTag, error) {
	// Get the first tag
	listTagsOptions := github.ListOptions{
		PerPage: 100,
	}
	_, response, err := c.baseClient.Repositories.ListTags(
		c.httpContext,
		c.RepoContext.Owner,
		c.RepoContext.Repo,
		&listTagsOptions,
	)
	if err != nil {
		return nil, err
	}

	listTagsOptions.Page = response.LastPage
	lastPageOfTags, _, err := c.baseClient.Repositories.ListTags(
		c.httpContext,
		c.RepoContext.Owner,
		c.RepoContext.Repo,
		&listTagsOptions,
	)
	if err != nil {
		return nil, err
	}

	return lastPageOfTags[len(lastPageOfTags)-1], nil
}

func (c *GitHubClient) GetLatestTag() (*github.RepositoryTag, error) {
	// Get the latest tag
	listTagsOptions := github.ListOptions{
		PerPage: 1,
	}
	tags, _, err := c.baseClient.Repositories.ListTags(
		c.httpContext,
		c.RepoContext.Owner,
		c.RepoContext.Repo,
		&listTagsOptions,
	)
	if err != nil {
		return nil, err
	}

	return tags[0], nil
}

func (c *GitHubClient) GetFirstCommit() (*github.Commit, error) {
	// Get the first commit
	listCommitsOptions := github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	_, response, err := c.baseClient.Repositories.ListCommits(
		c.httpContext,
		c.RepoContext.Owner,
		c.RepoContext.Repo,
		&listCommitsOptions,
	)
	if err != nil {
		return nil, err
	}

	listCommitsOptions.ListOptions.Page = response.LastPage
	lastPageOfCommits, _, err := c.baseClient.Repositories.ListCommits(
		c.httpContext,
		c.RepoContext.Owner,
		c.RepoContext.Repo,
		&listCommitsOptions,
	)
	if err != nil {
		return nil, err
	}

	firstCommit, _, err := c.baseClient.Git.GetCommit(
		c.httpContext,
		c.RepoContext.Owner,
		c.RepoContext.Repo,
		lastPageOfCommits[len(lastPageOfCommits)-1].GetSHA(),
	)
	if err != nil {
		return nil, err
	}

	return firstCommit, nil
}

func (c *GitHubClient) GetLatestCommit() (*github.Commit, error) {
	// Get the latest commit
	listCommitsOptions := github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	}

	commits, _, err := c.baseClient.Repositories.ListCommits(
		c.httpContext,
		c.RepoContext.Owner,
		c.RepoContext.Repo,
		&listCommitsOptions,
	)
	if err != nil {
		return nil, err
	}

	latestCommit, _, err := c.baseClient.Git.GetCommit(
		c.httpContext,
		c.RepoContext.Owner,
		c.RepoContext.Repo,
		commits[0].GetSHA(),
	)
	if err != nil {
		return nil, err
	}

	return latestCommit, nil
}
