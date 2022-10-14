// Package environment provides helper methods for working with the
// current environment.
package environment

import "os"

// IsCI returns true if the CI environment variable is set to true.
// This is used for most CI systems.
func IsCI() bool {
	return os.Getenv("CI") == "true"
}
