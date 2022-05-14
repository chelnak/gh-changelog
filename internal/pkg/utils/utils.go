//Package utils contains a number generic of methods that are used
//throughout the application.
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Masterminds/semver/v3"
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

func IsValidSemanticVersion(version string) bool {
	_, err := semver.NewVersion(version)
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

	if versionIsGreaterThan(currentVersion, release.Version) {
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

	body, err := ioutil.ReadAll(response.Body)
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

func versionIsGreaterThan(currentVersion, nextVersion string) bool {
	constraint, err := semver.NewConstraint(fmt.Sprintf("> %s", currentVersion))
	if err != nil {
		return false
	}

	version, err := semver.NewVersion(nextVersion)
	if err != nil {
		return false
	}

	return constraint.Check(version)
}

func parseLocalVersion(version string) string {
	slice := strings.Split(version, " ")

	if len(slice) == 1 {
		return version
	}

	return slice[2]
}
