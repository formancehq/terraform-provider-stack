package e2e_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestLedgerSchema(t *testing.T) {
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
							schema = {
								"test" = {
									"$segment1" = {
										".pattern" = "^[0-9]{10}$"
									}
								}
								"segment2" = {
									".metadata" = {
										"test" = {
											"default" = "test"
										}
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
