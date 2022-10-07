package changelog_test

import (
	"testing"
	"time"

	"github.com/chelnak/gh-changelog/internal/changelog"
	"github.com/chelnak/gh-changelog/internal/configuration"
	"github.com/chelnak/gh-changelog/internal/githubclient"
	"github.com/chelnak/gh-changelog/mocks"
	"github.com/stretchr/testify/assert"
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
	changelog.Now = func() time.Time {
		return safeParseTime()
	}

	mockGitHubClient.On("GetPullRequestsBetweenDates", safeParseTime(), changelog.Now()).Return([]githubclient.PullRequest{}, nil).Once()
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

func setupBuilder(gitClient *mocks.GitClient, gitHubClient *mocks.GitHubClient) changelog.ChangelogBuilder { // nolint:unparam
	_ = configuration.InitConfig()
	b := changelog.NewChangelogBuilder()

	if gitClient == nil {
		gitClient = setupMockGitClient()
	}

	if gitHubClient == nil {
		gitHubClient = setupMockGitHubClient()
	}

	b.WithSpinner(true)
	b.WithGitClient(gitClient)
	b.WithGitHubClient(gitHubClient)

	return b
}

func TestChangelogBuilder(t *testing.T) {
	builder := setupBuilder(nil, nil)

	changelog, err := builder.Build()
	assert.NoError(t, err)

	assert.Equal(t, "repo-name", changelog.GetRepoName())
	assert.Equal(t, "repo-owner", changelog.GetRepoOwner())

	assert.Len(t, changelog.GetUnreleased(), 0)
	assert.Len(t, changelog.GetEntries(), 2)

	assert.Equal(
		t,
		"this is a test pr 2 [#2](https://github.com/repo-owner/repo-name/pull/2) ([test-user](https://github.com/test-user))\n",
		changelog.GetEntries()[0].Added[0],
	)
}

func TestShouldErrorWithAnOlderNextVersion(t *testing.T) {
	builder := setupBuilder(nil, nil)
	builder.WithNextVersion("v0.0.1")

	_, err := builder.Build()
	assert.Error(t, err)
	assert.Equal(t, "the next version should be greater than the former: 'v0.0.1' ≤ 'v2.0.0'", err.Error())
}

func TestShouldErrorWithNoTags(t *testing.T) {
	mockGitHubClient := &mocks.GitHubClient{}
	mockGitHubClient.On("GetTags").Return([]githubclient.Tag{}, nil)

	builder := setupBuilder(nil, mockGitHubClient)

	_, err := builder.Build()
	assert.Error(t, err)
	assert.Equal(t, "there are no tags on this repository to evaluate", err.Error())
}

func TestWithFromVersion(t *testing.T) {
	builder := setupBuilder(nil, nil)
	builder.WithFromVersion("v2.0.0")

	changelog, err := builder.Build()
	assert.NoError(t, err)
	assert.Len(t, changelog.GetEntries(), 1)
	assert.Equal(
		t,
		"this is a test pr 2 [#2](https://github.com/repo-owner/repo-name/pull/2) ([test-user](https://github.com/test-user))\n",
		changelog.GetEntries()[0].Added[0],
	)
}

func TestWithFromLastVersion(t *testing.T) {
	builder := setupBuilder(nil, nil)
	builder.WithFromLastVersion(true)

	changelog, err := builder.Build()
	assert.NoError(t, err)
	assert.Len(t, changelog.GetEntries(), 1)
	assert.Equal(
		t,
		"this is a test pr 2 [#2](https://github.com/repo-owner/repo-name/pull/2) ([test-user](https://github.com/test-user))\n",
		changelog.GetEntries()[0].Added[0],
	)
}
