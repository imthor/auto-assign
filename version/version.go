package version

import "fmt"

var (
	// Version is the current version of the application
	Version = "0.1.0"
	// BuildTime is the time the binary was built
	BuildTime = "unknown"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
)

// String returns the version information as a string
func String() string {
	return fmt.Sprintf("Version: %s\nBuild Time: %s\nGit Commit: %s", Version, BuildTime, GitCommit)
}
