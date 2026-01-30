// Package version contains build-time version information.
package version

// Version is the current version of the application.
// It should be set at build time using ldflags:
//
//	-ldflags "-X github.com/BenbenIO/gh-quoi/internal/version.Version=v1.0.0"
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
