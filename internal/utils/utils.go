// Package utils contains a number generic of methods that are used
// throughout the application.
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/chelnak/gh-changelog/internal/version"
	"github.com/cli/go-gh"
	"github.com/fatih/color"
)

func SliceContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func IsValidSemanticVersion(v string) bool {
	_, err := version.NormalizeVersion(v)
	return err == nil
}

type Release struct {
	Version string `json:"tag_name"`
}

func CheckForUpdate(currentVersion string) bool {
	release, err := requestLatestRelease()
	if err != nil {
		return false
	}

	currentVersion = parseLocalVersion(currentVersion)

	if NextVersionIsGreaterThanCurrent(release.Version, currentVersion) {
		color.Yellow("\nVersion %s is available ✨\n\n", release.Version)
		fmt.Println("Run", color.GreenString(`gh extension upgrade chelnak/gh-changelog`), "to upgrade.")

		fmt.Println("\nAlternatively, you can disable this check by setting", color.GreenString("check_for_updates"), "to", color.RedString("false"), "via the configuration.")
		fmt.Println()
		return true
	}

	return false
}

func requestLatestRelease() (Release, error) {
	response, err := http.Get("https://api.github.com/repos/chelnak/gh-changelog/releases/latest")
	if err != nil {
		return Release{}, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Release{}, err
	}

	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return Release{}, err
	}

	return release, nil
}

func NextVersionIsGreaterThanCurrent(nextVersion, currentVersion string) bool {
	currentSemVer, err := version.NormalizeVersion(currentVersion)
	if err != nil {
		return false
	}

	// The nextVersion has already been validated by the builder
	// so we can safely eat the error.
	nextSemVer, err := version.NormalizeVersion(nextVersion)
	if err != nil {
		return false
	}

	return nextSemVer.GreaterThan(currentSemVer)
}

func parseLocalVersion(version string) string {
	slice := strings.Split(version, " ")

	if len(slice) == 1 {
		return version
	}

	return slice[2]
}

// RepoContext is a struct that contains the current repository owner and name.
type RepoContext struct {
	Owner string
	Name  string
}

// GetRepoContext returns a new RepoContext struct with the current repository owner and name.
func GetRepoContext() (RepoContext, error) {
	currentRepository, err := gh.CurrentRepository()
	if err != nil {
		if strings.Contains(err.Error(), "not a git repository (or any of the parent directories)") {
			return RepoContext{}, fmt.Errorf("the current directory is not a git repository")
		}

		return RepoContext{}, err
	}

	return RepoContext{
		Owner: currentRepository.Owner(),
		Name:  currentRepository.Name(),
	}, nil
}
