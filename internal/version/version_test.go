package version_test

import (
	"testing"

	"github.com/chelnak/gh-changelog/internal/version"
	"github.com/stretchr/testify/require"
)

func TestIsValidSemanticVersion(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "valid semantic version",
			value: "1.0.0",
			want:  true,
		},
		{
			name:  "valid semantic version with pre-release",
			value: "1.0.0-beta",
			want:  true,
		},
		{
			name:  "invalid semantic version",
			value: "asdasdasd1",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, version.IsValidSemanticVersion(tt.value))
		})
	}
}

func TestVersionIsGreaterThan(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "version is greater than",
			value: "2.0.0",
			want:  true,
		},
		{
			name:  "version is not greater than",
			value: "0.1.0",
			want:  false,
		},
		{
			name:  "when the version is equal",
			value: "1.0.0",
			want:  false,
		},
		{
			name:  "when the version is greater with pre-release",
			value: "1.0.1-beta",
			want:  true,
		},
		{
			name:  "version is not greater than with pre-release",
			value: "0.2.0-beta",
			want:  false,
		},
		{
			name:  "when the version is equal with pre-release",
			value: "1.0.0-beta",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, version.NextVersionIsGreaterThanCurrent(tt.value, "1.0.0"))
		})
	}
}

func TestVersionIsGreaterThanPreRelease(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "version is greater than and not a pre-release",
			value: "7.0.0",
			want:  true,
		},
		{
			name:  "version is greater than and is a standard release",
			value: "6.0.0",
			want:  true,
		},
		{
			name:  "version is greater than and is a pre-release",
			value: "6.0.1-rc.1",
			want:  true,
		},
		{
			name:  "version is not greater than and not a pre-repease",
			value: "0.1.0",
			want:  false,
		},
		{
			name:  "version is not greater than and is a pre-repease",
			value: "v0.1.0-rc.1",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, version.NextVersionIsGreaterThanCurrent(tt.value, "6.0.0-rc.1"))
		})
	}
}

func TestVersionParsesWithDifferentPreReleaseDelimeters(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "version is greater than and is a pre-release, using a -",
			value: "6.0.1-rc.1",
			want:  true,
		},
		{
			name:  "version is greater than and is a pre-release, using a .",
			value: "6.0.1-rc.1",
			want:  true,
		},
		{
			name:  "version is not greater than and is a pre-repease, using a -",
			value: "v0.1.0-rc.1",
			want:  false,
		},
		{
			name:  "version is not greater than and is a pre-repease, using a .",
			value: "v0.1.0.rc.1",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, version.NextVersionIsGreaterThanCurrent(tt.value, "6.0.0-rc.1"))
		})
	}
}
