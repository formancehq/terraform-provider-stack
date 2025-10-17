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

func TestLedgerDefault(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"cloud": providerserver.NewProtocol6WithError(CloudProvider()),
			"stack": providerserver.NewProtocol6WithError(StackProvider()),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			newTestStepStack(),
			{
				Config: newStack(RegionName) +
					`
						provider "stack" {
							stack_id = cloud_stack.default.id
							organization_id = data.cloud_current_organization.default.id
							uri = cloud_stack.default.uri
						}

						resource "stack_ledger" "default" {
							name = "test"
						}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("name"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("bucket"), knownvalue.StringExact("_default")),
				},
			},
			{
				Config: newStack(RegionName) +
					`
						provider "stack" {
							stack_id = cloud_stack.default.id
							organization_id = data.cloud_current_organization.default.id
							uri = cloud_stack.default.uri
						}

						resource "stack_ledger" "default" {
							name = "test"
						}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("name"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("bucket"), knownvalue.StringExact("_default")),
				},
			},
			{
				RefreshState: true,
			},
		},
	})

}

func TestLedgerWithMetadata(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"cloud": providerserver.NewProtocol6WithError(CloudProvider()),
			"stack": providerserver.NewProtocol6WithError(StackProvider()),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			newTestStepStack(),
			{
				Config: newStack(RegionName) +
					`
						provider "stack" {
							stack_id = cloud_stack.default.id
							organization_id = data.cloud_current_organization.default.id
							uri = cloud_stack.default.uri
						}

						resource "stack_ledger" "default" {
							name = "test"
							metadata = {
								"key1" = "value1"
								"key2" = "value2"
							}
						}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("name"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("bucket"), knownvalue.StringExact("_default")),
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("metadata"), knownvalue.MapExact(
						map[string]knownvalue.Check{
							"key1": knownvalue.StringExact("value1"),
							"key2": knownvalue.StringExact("value2"),
						},
					)),
				},
			},
			{
				Config: newStack(RegionName) +
					`
						provider "stack" {
							stack_id = cloud_stack.default.id
							organization_id = data.cloud_current_organization.default.id
							uri = cloud_stack.default.uri
						}

						resource "stack_ledger" "default" {
							name = "test"
							metadata = {
								"key1" = "value1"
								"key2" = "newvalue"
							}
						}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("name"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("bucket"), knownvalue.StringExact("_default")),
					statecheck.ExpectKnownValue("stack_ledger.default", tfjsonpath.New("metadata"), knownvalue.MapExact(
						map[string]knownvalue.Check{
							"key1": knownvalue.StringExact("value1"),
							"key2": knownvalue.StringExact("newvalue"),
						},
					)),
				},
			},
			{
				RefreshState: true,
			},
		},
	})

}
