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

func TestWebhooks(t *testing.T) {
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
				Config: newStack(OrganizationId, RegionName) +
					`
						resource "formancecloud_stack_module" "webhooks" {
							name = "webhooks"
							stack_id = formancecloud_stack.default.id
							organization_id = data.formancecloud_organizations.default.id
						}

						provider "formancestack" {
							stack_id = formancecloud_stack.default.id
							organization_id = formancecloud_stack.default.organization_id
							uri = formancecloud_stack.default.uri
						}

						resource "formancestack_webhooks" "webhooks" {
							endpoint = "https://formance.staging.com/webhook"
							event_types = [
								"transaction.created",
								"transaction.updated",
							]

							depends_on = [
								formancecloud_stack_module.webhooks,
							]
						}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("formancestack_webhooks.webhooks", tfjsonpath.New("id"), knownvalue.StringRegexp(regexp.MustCompile(`.+`))),
					statecheck.ExpectKnownValue("formancestack_webhooks.webhooks", tfjsonpath.New("endpoint"), knownvalue.StringExact("https://formance.staging.com/webhook")),
					statecheck.ExpectKnownValue("formancestack_webhooks.webhooks", tfjsonpath.New("event_types"), knownvalue.ListExact(
						[]knownvalue.Check{
							knownvalue.StringExact("transaction.created"),
							knownvalue.StringExact("transaction.updated"),
						},
					)),
					statecheck.ExpectSensitiveValue("formancestack_webhooks.webhooks", tfjsonpath.New("secret")),
				},
			},
			{
				Config: newStack(OrganizationId, RegionName) +
					`
							resource "formancecloud_stack_module" "webhooks" {
								name = "webhooks"
								stack_id = formancecloud_stack.default.id
								organization_id = data.formancecloud_organizations.default.id
							}
	
							provider "formancestack" {
								stack_id = formancecloud_stack.default.id
								organization_id = formancecloud_stack.default.organization_id
								uri = formancecloud_stack.default.uri
							}
	
							resource "formancestack_webhooks" "webhooks" {
								endpoint = "https://formance.staging.com/webhook2"
								event_types = [
									"transaction.test",
								]
	
								depends_on = [
									formancecloud_stack_module.webhooks,
								]
							}
						`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("formancestack_webhooks.webhooks", tfjsonpath.New("id"), knownvalue.StringRegexp(regexp.MustCompile(`.+`))),
					statecheck.ExpectKnownValue("formancestack_webhooks.webhooks", tfjsonpath.New("endpoint"), knownvalue.StringExact("https://formance.staging.com/webhook2")),
					statecheck.ExpectKnownValue("formancestack_webhooks.webhooks", tfjsonpath.New("event_types"), knownvalue.ListExact(
						[]knownvalue.Check{
							knownvalue.StringExact("transaction.test"),
						},
					)),
					statecheck.ExpectSensitiveValue("formancestack_webhooks.webhooks", tfjsonpath.New("secret")),
				},
			},
			{
				RefreshState: true,
			},
		},
	})

}
