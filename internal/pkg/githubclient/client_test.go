package githubclient_test

import (
	"testing"

	"github.com/chelnak/gh-changelog/internal/pkg/githubclient"
	"github.com/chelnak/gh-changelog/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_ItReturnsTheCorrectRepoName(t *testing.T) {
	mockClient := &mocks.GitHubClient{}
	mockClient.On("GetRepoName").Return("TestRepo")

	repoName := mockClient.GetRepoName()
	mockClient.AssertExpectations(t)

	assert.Equal(t, "TestRepo", repoName)
}

func Test_ItReturnsTheCorrectRepoOwner(t *testing.T) {
	mockClient := &mocks.GitHubClient{}
	mockClient.On("GetRepoOwner").Return("TestOwner")

	repoName := mockClient.GetRepoOwner()
	mockClient.AssertExpectations(t)

	assert.Equal(t, "TestOwner", repoName)
}

func NewJSONResponder(status int, body string) httpmock.Responder {
	resp := httpmock.NewStringResponse(status, body)
	resp.Header.Set("Content-Type", "application/json")
	return httpmock.ResponderFromResponse(resp)
}

func Test_GetTagsReturnsASliceOfTags(t *testing.T) {
	t.Skip() //Skip this test for now

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://api.github.com/graphql",
		NewJSONResponder(500, httpmock.File("data/get_tags_response.json").String()),
	)

	t.Setenv("GITHUB_TOKEN", "xxxxxxxx")
	t.Setenv("GH_TOKEN", "test-token")
	t.Setenv("GH_REPO", "test/repo")
	t.Setenv("GH_CONFIG_DIR", t.TempDir())

	client, err := githubclient.NewGitHubClient()

	assert.NoError(t, err)

	tags, err := client.GetTags()

	assert.NoError(t, err)

	assert.Equal(t, 2, len(tags))
	assert.Equal(t, "v1.0.0", tags[0].Name)
}
