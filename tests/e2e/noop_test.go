package e2e_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestNoopResources(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"formancecloud": providerserver.NewProtocol6WithError(CloudProvider()),
			"formancestack": providerserver.NewProtocol6WithError(StackProvider()),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			newTestStepStack(),
			{
				Config: newStack(RegionName) +
					`
						provider "formancestack" {
							stack_id = formancecloud_stack.default.id
							organization_id = data.formancecloud_current_organization.default.id
							uri = formancecloud_stack.default.uri
						}

						resource "formancestack_noop" "noop" {}
					
					`,
			},
		},
	})

}
