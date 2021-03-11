package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

type Info struct {
	Major        string `json:"major"`
	Minor        string `json:"minor"`
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

func (i Info) Format() (string, error) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

var (
	GitVersion   = "0.0.0-dev"
	gitCommit    string
	gitTreeState string
	buildDate    = "1970-01-01T00:00:00Z"
)

// These variables come from -ldflags
func Get() Info {
	var (
		version  = strings.Split(GitVersion, ".")
		gitMajor string
		gitMinor string
	)

	if len(version) >= 2 {
		gitMajor = version[0]
		gitMinor = version[1]
	}

	return Info{
		Major:        gitMajor,
		Minor:        gitMinor,
		GitVersion:   GitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
