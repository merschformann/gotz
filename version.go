package main

// Version, commit, and date are set by the release process.
var (
	version = "v0.0.0"
	commit  = "none"
	date    = "1970-01-01T00:00:00Z"
)

type ReleaseInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// GetReleaseInfo returns the version information.
func GetReleaseInfo() ReleaseInfo {
	return ReleaseInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}
