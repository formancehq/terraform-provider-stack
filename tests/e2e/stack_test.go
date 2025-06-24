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

func TestStack(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"formancecloud": providerserver.NewProtocol6WithError(CloudProvider()),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			{
				Config: newStack(RegionName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("formancecloud_stack.default", tfjsonpath.New("name"), knownvalue.StringExact("test")),
					statecheck.ExpectKnownValue("formancecloud_stack.default", tfjsonpath.New("id"), knownvalue.StringRegexp(regexp.MustCompile(`.+`))),
					statecheck.ExpectKnownValue("formancecloud_stack.default", tfjsonpath.New("force_destroy"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("formancecloud_stack.default", tfjsonpath.New("uri"), knownvalue.StringRegexp(regexp.MustCompile(`.+`))),
				},
			},
		},
	})

}
