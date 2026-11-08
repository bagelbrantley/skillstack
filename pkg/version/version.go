package version

import "fmt"

var (
	Version   = "v1.0.0"
	GitCommit = "dev"
	BuildDate = "unknown"
)

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
