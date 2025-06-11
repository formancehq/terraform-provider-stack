package internal

import "fmt"

const (
	Repository  = "formancehq/terraform-provider-cloud"
	ServiceName = "terraform-provider-cloud"
)

type AppInfo struct {
	Name       string
	Repository string
	Version    string
	BuildDate  string
	Commit     string
}

func (a AppInfo) String() string {
	return fmt.Sprintf("\n\tName: %s\n\tVersion: %s\n\tBuildDate: %s\n\tCommit: %s\n\tRepository: %s\n\t", a.Name, a.Version, a.BuildDate, a.Commit, a.Repository)
}
