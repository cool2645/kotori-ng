package version

var (
	BuildTag      string
	BuildTime     string
	GitCommitSHA1 string
	GitTag        string
	vi            = VersionInfo{
		Name:    "kotori-ng",
		Author:  "2645 Studio",
		Version: "1.0-pre-alpha.1",
		License: "Unlicense",
		URL:     "https://github.com/cool2645/kotori-ng",
	}
	bi = BuildInfo{
		BuildTag:      BuildTag,
		BuildTime:     BuildTime,
		GitCommitSHA1: GitCommitSHA1,
		GitTag:        GitTag,
	}
)

type VersionInfo struct {
	Name    string `json:"name"`
	Author  string `json:"author"`
	Version string `json:"version"`
	License string `json:"license"`
	URL     string `json:"url"`
}

type BuildInfo struct {
	BuildTag      string `json:"build_tag"`
	BuildTime     string `json:"build_time"`
	GitCommitSHA1 string `json:"git_commit_sha1"`
	GitTag        string `json:"git_tag"`
}

func GetBuildInfo() BuildInfo {
	return bi
}

func GetVersionInfo() VersionInfo {
	return vi
}
