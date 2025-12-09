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

func TestPaymentPoolQuery(t *testing.T) {
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

						resource "stack_payments_pool" "default" {
							name = "default_pool"
							query = {
								"$match" = {
									"id" = "1234567890"
								}
							}
						}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("stack_payments_pool.default", tfjsonpath.New("query"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"$match": knownvalue.ObjectExact(map[string]knownvalue.Check{
							"id": knownvalue.StringExact("1234567890"),
						}),
					})),
				},
			},
			{
				RefreshState: true,
			},
		},
	})

}
