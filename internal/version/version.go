// Package version contains a wrapper for parsing semantic versions with the
// semver library.
// The code here will be removed if/when semver can handle pre-release versions with a dot.
// It's not ideal but a reasonable workaround for now.
package version

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// The compiled version of the regex created at init() is cached here so it
// only needs to be created once.
var versionRegex *regexp.Regexp

// semVerRegex is the regular expression used to parse a semantic version.
const semVerRegex string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`((?:-|\.)([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

// Version represents a single semantic version.
type Version struct {
	major, minor, patch uint64
	pre                 string
	metadata            string
	original            string
}

func init() {
	versionRegex = regexp.MustCompile("^" + semVerRegex + "$")
}

const (
	num     string = "0123456789"
	allowed string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-" + num
)

// NormalizeVersion is basically a copy of SemVers NewVersion. However, it's
// purpose is to normalize the version string to a semver compatible one.
func NormalizeVersion(v string) (*semver.Version, error) {
	m := versionRegex.FindStringSubmatch(v)
	if m == nil {
		return nil, semver.ErrInvalidSemVer
	}

	sv := &Version{
		metadata: m[8],
		pre:      m[5],
		original: v,
	}

	var err error
	sv.major, err = strconv.ParseUint(m[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing version segment: %s", err)
	}

	if m[2] != "" {
		sv.minor, err = strconv.ParseUint(strings.TrimPrefix(m[2], "."), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing version segment: %s", err)
		}
	} else {
		sv.minor = 0
	}

	if m[3] != "" {
		sv.patch, err = strconv.ParseUint(strings.TrimPrefix(m[3], "."), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing version segment: %s", err)
		}
	} else {
		sv.patch = 0
	}

	// Perform some basic due diligence on the extra parts to ensure they are
	// valid.

	if sv.pre != "" {
		if err = validatePrerelease(sv.pre); err != nil {
			return nil, err
		}
	}

	if sv.metadata != "" {
		if err = validateMetadata(sv.metadata); err != nil {
			return nil, err
		}
	}

	// Return the semver version.
	return semver.NewVersion(sv.String())
}

// String converts a Version object to a string.
// Note, if the original version contained a leading v this version will not.
// See the Original() method to retrieve the original value. Semantic Versions
// don't contain a leading v per the spec. Instead it's optional on
// implementation.
func (v Version) String() string {
	var buf bytes.Buffer

	_, _ = fmt.Fprintf(&buf, "%d.%d.%d", v.major, v.minor, v.patch)
	if v.pre != "" {
		_, _ = fmt.Fprintf(&buf, "-%s", v.pre)
	}
	if v.metadata != "" {
		_, _ = fmt.Fprintf(&buf, "+%s", v.metadata)
	}

	return buf.String()
}

// Like strings.ContainsAny but does an only instead of any.
func containsOnly(s string, comp string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(comp, r)
	}) == -1
}

// From the spec, "Identifiers MUST comprise only
// ASCII alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty.
// Numeric identifiers MUST NOT include leading zeroes.". These segments can
// be dot separated.
func validatePrerelease(p string) error {
	eparts := strings.Split(p, ".")
	for _, p := range eparts {
		if containsOnly(p, num) {
			if len(p) > 1 && p[0] == '0' {
				return semver.ErrSegmentStartsZero
			}
		} else if !containsOnly(p, allowed) {
			return semver.ErrInvalidPrerelease
		}
	}

	return nil
}

// From the spec, "Build metadata MAY be denoted by
// appending a plus sign and a series of dot separated identifiers immediately
// following the patch or pre-release version. Identifiers MUST comprise only
// ASCII alphanumerics and hyphen [0-9A-Za-z-]. Identifiers MUST NOT be empty."
func validateMetadata(m string) error {
	eparts := strings.Split(m, ".")
	for _, p := range eparts {
		if !containsOnly(p, allowed) {
			return semver.ErrInvalidMetadata
		}
	}
	return nil
}

func IsValidSemanticVersion(v string) bool {
	_, err := NormalizeVersion(v)
	return err == nil
}

func NextVersionIsGreaterThanCurrent(nextVersion, currentVersion string) bool {
	currentSemVer, err := NormalizeVersion(currentVersion)
	if err != nil {
		return false
	}

	// The nextVersion has already been validated by the builder
	// so we can safely eat the error.
	nextSemVer, err := NormalizeVersion(nextVersion)
	if err != nil {
		return false
	}

	return nextSemVer.GreaterThan(currentSemVer)
}
