package internal

import "fmt"

const (
	TerraformRepository = "formancehq/stack"
	ServiceName         = "terraform-provider-stack"
)

type AppInfo struct {
	Name                string
	TerraformRepository string
	Version             string
	BuildDate           string
	Commit              string
}

func (a AppInfo) String() string {
	return fmt.Sprintf("\n\tName: %s\n\tVersion: %s\n\tBuildDate: %s\n\tCommit: %s\n\tTerraform Repository: %s\n\t", a.Name, a.Version, a.BuildDate, a.Commit, a.TerraformRepository)
}
