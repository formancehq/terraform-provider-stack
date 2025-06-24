// Format Terraform code for use in documentation.
// If you do not have Terraform installed, you can remove the formatting command, but it is suggested
// to ensure the documentation is formatted properly.

//go:generate rm -rf docs
//go:generate mkdir docs
//go:generate terraform fmt -recursive ./examples/

// Generate documentation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest generate --provider-dir . -provider-name formancestack

package main

import "github.com/formancehq/terraform-provider-stack/cmd"

func main() {
	cmd.Execute()
}
