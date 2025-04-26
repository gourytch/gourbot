package version

import "fmt"

var (
	Version    = "n/a"
	CommitHash = "n/a"
	BuildTime  = "n/a"
)

func VersionString() string {
	return fmt.Sprintf("version: %s, CommitHash: %s, BuildTime: %s", Version, CommitHash, BuildTime)
}
