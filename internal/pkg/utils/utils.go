package utils

import (
	"github.com/Masterminds/semver/v3"
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
