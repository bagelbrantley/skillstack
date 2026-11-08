package version

import (
	"fmt"
	"runtime/debug"
)

var (
	Version   = "v1.0.1"
	GitCommit = "dev"
	BuildDate = "unknown"
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if len(s.Value) >= 7 {
					GitCommit = s.Value[:7]
				} else {
					GitCommit = s.Value
				}
			case "vcs.time":
				BuildDate = s.Value
			}
		}
	}
}

type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
}

func GetInfo() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
	}
}

func (i Info) String() string {
	return fmt.Sprintf("skillstack %s (commit: %s, built: %s)", i.Version, i.GitCommit, i.BuildDate)
}
