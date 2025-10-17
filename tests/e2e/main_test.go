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
	"go.opentelemetry.io/otel"
)

var (
	CloudProvider func() provider.Provider
	StackProvider func() provider.Provider
	RegionName    = ""
)

func newTestStepStack() resource.TestStep {
	return resource.TestStep{
		Config: newStack(RegionName),
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue("cloud_stack.default", tfjsonpath.New("name"), knownvalue.StringExact("test")),
			statecheck.ExpectKnownValue("cloud_stack.default", tfjsonpath.New("id"), knownvalue.StringRegexp(regexp.MustCompile(`.+`))),
			statecheck.ExpectKnownValue("cloud_stack.default", tfjsonpath.New("force_destroy"), knownvalue.Bool(true)),
			statecheck.ExpectKnownValue("cloud_stack.default", tfjsonpath.New("uri"), knownvalue.StringRegexp(regexp.MustCompile(`.+`))),
		},
	}
}

func TestMain(m *testing.M) {
	endpoint := os.Getenv("FORMANCE_CLOUD_API_ENDPOINT")
	clientID := os.Getenv("FORMANCE_CLOUD_CLIENT_ID")
	clientSecret := os.Getenv("FORMANCE_CLOUD_CLIENT_SECRET")

	// TODO: This can be replaced with TF variables in the future
	RegionName = os.Getenv("FORMANCE_CLOUD_REGION_NAME")
	if RegionName == "" || endpoint == "" || clientID == "" || clientSecret == "" {
		missingVars := []string{RegionName, endpoint, clientID, clientSecret}
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
		otel.GetTracerProvider(),
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
		otel.GetTracerProvider(),
		logging.Testing(),
		endpoint,
		clientID,
		clientSecret,
		transport,
		cloudpkg.NewCloudSDK(),
		cloudpkg.NewTokenProvider,
	)

	code := m.Run()

	os.Exit(code)
}

func newStack(regionName string) string {
	return `
		data "cloud_current_organization" "default" {}

		data "cloud_regions" "default" {
			name = "` + regionName + `"
		}

		resource "cloud_stack" "default" {
			name = "test"
			region_id = data.cloud_regions.default.id

			force_destroy = true
		}
	`
}
