package changelog_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/chelnak/gh-changelog/internal/pkg/changelog"
	"github.com/chelnak/gh-changelog/internal/pkg/configuration"
	"github.com/chelnak/gh-changelog/internal/pkg/githubclient"
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
			Name: "v1.0.0",
			Sha:  "42d4c93b23eaf307c5f9712f4c62014fe38332bd",
			Date: safeParseTime(),
		},
		{
			Name: "v2.0.0",
			Sha:  "0d724ba5b4235aa88d45a20f4ecd8db4b4695cf1",
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
			Number: 1,
			Title:  "this is a test pr",
			User:   "test-user",
			Labels: []githubclient.Label{
				{
					Name: "enhancement",
				},
			},
		},
		{
			Number: 2,
			Title:  "this is a test pr 2",
			User:   "test-user",
			Labels: []githubclient.Label{
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

var testBuilder = changelog.NewChangelogBuilder()

func TestChangelogBuilder(t *testing.T) {
	_ = configuration.InitConfig()

	mockGitClient := setupMockGitClient()
	mockGitHubClient := setupMockGitHubClient()

	testBuilder.WithSpinner(true)
	testBuilder.WithGitClient(mockGitClient)
	testBuilder.WithGitHubClient(mockGitHubClient)

	changelog, err := testBuilder.Build()
	assert.NoError(t, err)

	assert.Equal(t, changelog.GetRepoName(), "repo-name")
	assert.Equal(t, changelog.GetRepoOwner(), "repo-owner")

	assert.Len(t, changelog.GetUnreleased(), 0)
	assert.Len(t, changelog.GetEntries(), 2)

	fmt.Println(changelog.GetEntries())
	assert.Equal(t, changelog.GetEntries()[0].Added[0], "this is a test pr [#1](https://github.com/repo-owner/repo-name/pull/1) ([test-user](https://github.com/test-user))\n")
}
