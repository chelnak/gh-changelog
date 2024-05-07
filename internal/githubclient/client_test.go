package githubclient_test

import (
	"testing"

	mocks "github.com/chelnak/gh-changelog/mocks/githubclient"
	"github.com/stretchr/testify/require"
)

func Test_ItReturnsTheCorrectRepoName(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{}
	mockClient.EXPECT().GetRepoName().Return("TestRepo")
	repoName := mockClient.GetRepoName()
	require.Equal(t, "TestRepo", repoName)
}

func Test_ItReturnsTheCorrectRepoOwner(t *testing.T) {
	mockClient := &mocks.MockGitHubClient{}
	mockClient.EXPECT().GetRepoOwner().Return("TestOwner")
	repoName := mockClient.GetRepoOwner()
	require.Equal(t, "TestOwner", repoName)
}
