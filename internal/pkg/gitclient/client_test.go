package gitclient_test

import (
	"testing"

	"github.com/chelnak/gh-changelog/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_ItReturnsTheCorrectRepoName(t *testing.T) {
	mockClient := &mocks.GitHubClient{}
	mockClient.On("GetRepoName").Return("TestRepo")

	repoName := mockClient.GetRepoName()
	mockClient.AssertExpectations(t)

	assert.Equal(t, "TestRepo", repoName)
}

func TestItReturnsTheCorrectRepoOwner(t *testing.T) {
	mockClient := &mocks.GitHubClient{}
	mockClient.On("GetRepoOwner").Return("TestOwner")

	repoName := mockClient.GetRepoOwner()
	mockClient.AssertExpectations(t)

	assert.Equal(t, "TestOwner", repoName)
}
