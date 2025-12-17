package e2e_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestLedgerSchemaChart(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"cloud": providerserver.NewProtocol6WithError(CloudProvider()),
			"stack": providerserver.NewProtocol6WithError(StackProvider()),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
						provider "stack" {
							stack_id = cloud_stack.default.id
							organization_id = data.cloud_current_organization.default.id
							uri = cloud_stack.default.uri
						}

						data "cloud_current_organization" "default" {}

						data "cloud_regions" "default" {
							name = "` + RegionName + `"
						}

						resource "cloud_stack" "default" {
							name = "test"
							region_id = data.cloud_regions.default.id
							version = "v3.2-rc.1"
							force_destroy = true
						}

						resource "stack_ledger" "default" {
							name = "test"
						}

						resource "stack_ledger_schema" "default" {
							ledger = stack_ledger.default.name
							version = "v1.0.0"
							chart = {
								"banks"={
									"$iban"={
										".pattern"="DE*"
										main=null
										
									}
								}
							}
						}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("name"), knownvalue.StringExact("test")),
				},
			},
			{
				RefreshState: true,
			},
		},
	})

}

func TestLedgerSchemaTransactions(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"cloud": providerserver.NewProtocol6WithError(CloudProvider()),
			"stack": providerserver.NewProtocol6WithError(StackProvider()),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
						provider "stack" {
							stack_id = cloud_stack.default.id
							organization_id = data.cloud_current_organization.default.id
							uri = cloud_stack.default.uri
						}

						data "cloud_current_organization" "default" {}

						data "cloud_regions" "default" {
							name = "` + RegionName + `"
						}

						resource "cloud_stack" "default" {
							name = "test"
							region_id = data.cloud_regions.default.id
							version = "v3.2-rc.1"
							force_destroy = true
						}

						resource "stack_ledger" "default" {
							name = "test"
						}

						resource "stack_ledger_schema" "default" {
							ledger = stack_ledger.default.name
							version = "v1.0.0"
							transactions = {
								customer_deposit = {
									description = "Test transaction"
									script = <<-EOT
										vars {
											account $b
										}
										send [USD 100] (
											source = @world
											destination = $b
										)
									EOT
								}
							}
						}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("name"), knownvalue.StringExact("test")),
				},
			},
			{
				RefreshState: true,
			},
		},
	})

}

// TOfix(aini encoding ??): invalid template customer_deposit: compilation error: \\u001b[31m--\\u003e\\u001b[0m error:1:6\\n  \\u001b[34m|\\u001b[0m\\n\\u001b[31m1 | \\u001b[0m\\u001b[90mvars {\\u001b[0m}\\u001b[90m\\n\\u001b[0m  \\u001b[34m|\\u001b[0m       \\u001b[31m^\\u001b[0m extraneous input '}' expecting NEWLINE\\n\\u001b[31m--\\u003e\\u001b[0m error:2:0\\n  \\u001b[34m|\\u001b[0m\\n\\u001b[31m2 | \\u001b[0m\\u001b[90m\\u001b[0msend\\u001b[90m [USD 100] (\\n\\u001b[0m  \\u001b[34m|\\u001b[0m \
// \u001b[31m^^^\\u001b[0m mismatched input 'send' expecting {'account', 'asset', 'number', 'monetary', 'portion', 'string'}\\n\"}\n
func TestLedgerSchemaTransactionsScriptError(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"cloud": providerserver.NewProtocol6WithError(CloudProvider()),
			"stack": providerserver.NewProtocol6WithError(StackProvider()),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
						provider "stack" {
							stack_id = cloud_stack.default.id
							organization_id = data.cloud_current_organization.default.id
							uri = cloud_stack.default.uri
						}

						data "cloud_current_organization" "default" {}

						data "cloud_regions" "default" {
							name = "` + RegionName + `"
						}

						resource "cloud_stack" "default" {
							name = "test"
							region_id = data.cloud_regions.default.id
							version = "v3.2-rc.1"
							force_destroy = true
						}

						resource "stack_ledger" "default" {
							name = "test"
						}

						resource "stack_ledger_schema" "default" {
							ledger = stack_ledger.default.name
							version = "v1.0.0"
							transactions = {
								customer_deposit = {
									description = "Test transaction"
									script = <<-EOT
										vars {}
										send [USD 100] (
											source = @world
											destination = $b
										)
									EOT
								}
							}
						}
					`,
				ExpectError: regexp.MustCompile(""),
			},
		},
	})

}
