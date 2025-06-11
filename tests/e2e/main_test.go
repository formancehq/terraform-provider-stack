package e2e_test

import (
	"os"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/pkg/testprovider"
	"github.com/formancehq/terraform-provider-stack/internal/server"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

var (
	CloudProvider  func() provider.Provider
	StackProvider  func() provider.Provider
	RegionName     = ""
	OrganizationId = ""
)

func TestMain(m *testing.M) {
	endpoint := os.Getenv("FORMANCE_CLOUD_API_ENDPOINT")
	clientID := os.Getenv("FORMANCE_CLOUD_CLIENT_ID")
	clientSecret := os.Getenv("FORMANCE_CLOUD_CLIENT_SECRET")

	// This can be replaced with TF variables in the future
	RegionName = os.Getenv("FORMANCE_CLOUD_REGION_NAME")
	OrganizationId = os.Getenv("FORMANCE_CLOUD_ORGANIZATION_ID")

	StackProvider = server.NewStackProvider(
		logging.Testing(),
		server.FormanceStackEndpoint(endpoint),
		server.FormanceStackClientId(clientID),
		server.FormanceStackClientSecret(clientSecret),
		pkg.NewSDK,
	)
	CloudProvider = testprovider.NewCloudProvider(
		logging.Testing(),
		"develop",
		endpoint,
		clientID,
		clientSecret,
		pkg.NewSDK,
	)

	code := m.Run()

	os.Exit(code)
}

func newStack(organizationId, regionName string) string {
	return `
		data "formancecloud_organizations" "default" {
			id = "` + organizationId + `"
		}

		data "formancecloud_regions" "default" {
			name = "` + regionName + `"
			organization_id = data.formancecloud_organizations.default.id
		}

		resource "formancecloud_stack" "default" {
			name = "test"
			organization_id = data.formancecloud_organizations.default.id
			region_id = data.formancecloud_regions.default.id

			force_destroy = true
		}
	`
}
