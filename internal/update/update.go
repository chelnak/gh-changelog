// Package update provides functionality to check for updates to the extension.
package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/chelnak/gh-changelog/internal/version"
	"github.com/fatih/color"
)

type Release struct {
	Version string `json:"tag_name"`
}

func CheckForUpdate(currentVersion string) bool {
	release, err := requestLatestRelease()
	if err != nil {
		return false
	}

	currentVersion = parseLocalVersion(currentVersion)

	if version.NextVersionIsGreaterThanCurrent(release.Version, currentVersion) {
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

func parseLocalVersion(version string) string {
	slice := strings.Split(version, " ")

	if len(slice) == 1 {
		return version
	}

	return slice[2]
}
