package e2e_test

import (
	"flag"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/httpclient"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/otlp"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/pkg/testprovider"
	"github.com/formancehq/terraform-provider-stack/internal/server"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var (
	CloudProvider  func() provider.Provider
	StackProvider  func() provider.Provider
	RegionName     = ""
	OrganizationId = ""
)

func newTestStepStack() resource.TestStep {
	return resource.TestStep{
		Config: newStack(OrganizationId, RegionName),
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue("formancecloud_stack.default", tfjsonpath.New("name"), knownvalue.StringExact("test")),
			statecheck.ExpectKnownValue("formancecloud_stack.default", tfjsonpath.New("id"), knownvalue.StringRegexp(regexp.MustCompile(`.+`))),
			statecheck.ExpectKnownValue("formancecloud_stack.default", tfjsonpath.New("force_destroy"), knownvalue.Bool(true)),
			statecheck.ExpectKnownValue("formancecloud_stack.default", tfjsonpath.New("uri"), knownvalue.StringRegexp(regexp.MustCompile(`.+`))),
		},
	}
}

func TestMain(m *testing.M) {
	endpoint := os.Getenv("FORMANCE_CLOUD_API_ENDPOINT")
	clientID := os.Getenv("FORMANCE_CLOUD_CLIENT_ID")
	clientSecret := os.Getenv("FORMANCE_CLOUD_CLIENT_SECRET")

	// TODO: This can be replaced with TF variables in the future
	RegionName = os.Getenv("FORMANCE_CLOUD_REGION_NAME")
	OrganizationId = os.Getenv("FORMANCE_CLOUD_ORGANIZATION_ID")
	if RegionName == "" || OrganizationId == "" || endpoint == "" || clientID == "" || clientSecret == "" {
		missingVars := []string{RegionName, OrganizationId, endpoint, clientID, clientSecret}
		missingVars = collectionutils.Filter(missingVars, func(s string) bool {
			return s == ""
		})
		if len(missingVars) > 0 {
			panic("You must set the following environment variables: " + strings.Join(missingVars, ", "))
		}
	}

	flag.Parse()

	var transport http.RoundTripper
	if testing.Verbose() {
		transport = httpclient.NewDebugHTTPTransport(
			otlp.NewRoundTripper(http.DefaultTransport, true),
		)
	} else {
		transport = otlp.NewRoundTripper(http.DefaultTransport, false)
	}

	StackProvider = server.NewStackProvider(
		logging.Testing(),
		server.FormanceStackEndpoint(endpoint),
		server.FormanceStackClientId(clientID),
		server.FormanceStackClientSecret(clientSecret),
		transport,
		sdk.NewCloudSDK(),
		cloudpkg.NewTokenProvider,
		pkg.NewTokenProviderFn(),
		sdk.NewStackSdk(),
	)
	CloudProvider = testprovider.NewCloudProvider(
		logging.Testing(),
		endpoint,
		clientID,
		clientSecret,
		transport,
		cloudpkg.NewSDK,
		cloudpkg.NewTokenProvider,
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
