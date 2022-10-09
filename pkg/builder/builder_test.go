package builder_test

import (
	"testing"
	"time"

	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/githubclient"
	"github.com/chelnak/gh-changelog/mocks"
	"github.com/chelnak/gh-changelog/pkg/builder"
	"github.com/stretchr/testify/assert"
)

const (
	repoName  = "repo-name"
	repoOwner = "repo-owner"
)

func safeParseTime() time.Time {
	time, _ := time.Parse(time.RFC3339, time.Time{}.String())
	return time
}

func setupMockGitClient() *mocks.GitClient {
	mockGitClient := &mocks.GitClient{}
	mockGitClient.On("GetFirstCommit").Return("42d4c93b23eaf307c5f9712f4c62014fe38332bd", nil)
	mockGitClient.On("GetLastCommit").Return("0d724ba5b4235aa88d45a20f4ecd8db4b4695cf1", nil)
	mockGitClient.On("GetDateOfHash", "42d4c93b23eaf307c5f9712f4c62014fe38332bd").Return(safeParseTime(), nil).Once()
	return mockGitClient
}

func setupMockGitHubClient() *mocks.GitHubClient {
	mockGitHubClient := &mocks.GitHubClient{}
	mockGitHubClient.On("GetTags").Return([]githubclient.Tag{
		{
			Name: "v2.0.0",
			Sha:  "0d724ba5b4235aa88d45a20f4ecd8db4b4695cf1",
			Date: safeParseTime(),
		},
		{
			Name: "v1.0.0",
			Sha:  "42d4c93b23eaf307c5f9712f4c62014fe38332bd",
			Date: safeParseTime(),
		},
	}, nil)

	// bad ??
	builder.Now = func() time.Time {
		return safeParseTime()
	}

	mockGitHubClient.On("GetPullRequestsBetweenDates", safeParseTime(), builder.Now()).Return([]githubclient.PullRequest{}, nil).Once()
	mockGitHubClient.On("GetPullRequestsBetweenDates", time.Time{}, time.Time{}).Return([]githubclient.PullRequest{
		{
			Number: 2,
			Title:  "this is a test pr 2",
			User:   "test-user",
			Labels: []githubclient.PullRequestLabel{
				{
					Name: "enhancement",
				},
			},
		},
		{
			Number: 1,
			Title:  "this is a test pr",
			User:   "test-user",
			Labels: []githubclient.PullRequestLabel{
				{
					Name: "enhancement",
				},
			},
		},
	}, nil)

	mockGitHubClient.On("GetRepoName").Return(repoName)
	mockGitHubClient.On("GetRepoOwner").Return(repoOwner)

	return mockGitHubClient
}

func setupBuilder(opts *builder.BuilderOptions) builder.Builder { // nolint:unparam
	_ = configuration.InitConfig()

	if opts == nil {
		opts = &builder.BuilderOptions{}
	}

	opts.EnableSpinner = true

	if opts.GitClient == nil {
		opts.GitClient = setupMockGitClient()
	}

	if opts.GitHubClient == nil {
		opts.GitHubClient = setupMockGitHubClient()
	}

	b, _ := builder.NewBuilder(*opts)

	return b
}

func TestChangelogBuilder(t *testing.T) {
	builder := setupBuilder(nil)

	changelog, err := builder.BuildChangelog()
	assert.NoError(t, err)

	assert.Equal(t, repoName, changelog.GetRepoName())
	assert.Equal(t, repoOwner, changelog.GetRepoOwner())

	assert.Len(t, changelog.GetUnreleased(), 0)
	assert.Len(t, changelog.GetEntries(), 2)

	assert.Equal(
		t,
		"this is a test pr 2 [#2](https://github.com/repo-owner/repo-name/pull/2) ([test-user](https://github.com/test-user))\n",
		changelog.GetEntries()[0].Added[0],
	)
}

func TestShouldErrorWithAnOlderNextVersion(t *testing.T) {
	opts := &builder.BuilderOptions{
		NextVersion: "v0.0.1",
	}
	builder := setupBuilder(opts)
	_, err := builder.BuildChangelog()

	assert.Error(t, err)
	assert.Equal(t, "the next version should be greater than the former: 'v0.0.1' ≤ 'v2.0.0'", err.Error())
}

func TestShouldErrorWithNoTags(t *testing.T) {
	mockGitHubClient := &mocks.GitHubClient{}
	mockGitHubClient.On("GetTags").Return([]githubclient.Tag{}, nil)
	mockGitHubClient.On("GetRepoName").Return(repoName)
	mockGitHubClient.On("GetRepoOwner").Return(repoOwner)

	opts := &builder.BuilderOptions{
		GitHubClient: mockGitHubClient,
	}

	builder := setupBuilder(opts)
	_, err := builder.BuildChangelog()

	assert.Error(t, err)
	assert.Equal(t, "there are no tags on this repository to evaluate", err.Error())
}

func TestWithFromVersion(t *testing.T) {
	opts := &builder.BuilderOptions{
		FromVersion: "v2.0.0",
	}

	builder := setupBuilder(opts)
	changelog, err := builder.BuildChangelog()

	assert.NoError(t, err)
	assert.Len(t, changelog.GetEntries(), 1)
	assert.Equal(
		t,
		"this is a test pr 2 [#2](https://github.com/repo-owner/repo-name/pull/2) ([test-user](https://github.com/test-user))\n",
		changelog.GetEntries()[0].Added[0],
	)
}

func TestWithFromLastVersion(t *testing.T) {
	opts := &builder.BuilderOptions{
		LatestVersion: true,
	}

	builder := setupBuilder(opts)
	changelog, err := builder.BuildChangelog()

	assert.NoError(t, err)
	assert.Len(t, changelog.GetEntries(), 1)
	assert.Equal(
		t,
		"this is a test pr 2 [#2](https://github.com/repo-owner/repo-name/pull/2) ([test-user](https://github.com/test-user))\n",
		changelog.GetEntries()[0].Added[0],
	)
}
