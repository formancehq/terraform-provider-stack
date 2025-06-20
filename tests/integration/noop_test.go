package integration_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/go-libs/v3/logging"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-stack/internal/server"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"go.uber.org/mock/gomock"
)

func TestNoop(t *testing.T) {
	t.Parallel()

	type testCase struct {
		clientId     string
		clientSecret string
		endpoint     string

		organizationId string
		stackId        string
		uri            string

		expectedError string
	}

	for _, tc := range []testCase{
		{
			clientId:       fmt.Sprintf("organization_%s", uuid.NewString()[:8]),
			clientSecret:   uuid.NewString(),
			organizationId: uuid.NewString(),
			stackId:        uuid.NewString(),
			uri:            fmt.Sprintf("https://%s-%s.formance.cloud/api", uuid.NewString()[:8], uuid.NewString()[:4]),
		},
		{
			clientId:       uuid.NewString(),
			clientSecret:   uuid.NewString(),
			organizationId: uuid.NewString(),
			stackId:        uuid.NewString(),
			uri:            fmt.Sprintf("https://%s-%s.formance.cloud/api", uuid.NewString()[:8], uuid.NewString()[:4]),
			expectedError:  "Invalid client_id",
		}, {
			clientId: uuid.NewString(),

			organizationId: uuid.NewString(),
			stackId:        uuid.NewString(),
			uri:            fmt.Sprintf("https://%s-%s.formance.cloud/api", uuid.NewString()[:8], uuid.NewString()[:4]),
			expectedError:  "Missing client_secret",
		},
	} {
		t.Run(fmt.Sprintf("%s %+v", t.Name(), tc), func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			cloudSdk := sdk.NewMockCloudSDK(ctrl)
			tokenProvider, _ := cloudpkg.NewMockTokenProvider(ctrl)
			stackTokenProvider := pkg.NewMockTokenProviderImpl(ctrl)
			stacksdk := sdk.NewMockStackSdkImpl(ctrl)

			stackProvider := server.NewStackProvider(
				logging.Testing().WithField("provider", "stack_noop"),
				server.FormanceStackEndpoint(tc.endpoint),
				server.FormanceStackClientId(tc.clientId),
				server.FormanceStackClientSecret(tc.clientSecret),
				transport,
				func(creds cloudpkg.Creds, transport http.RoundTripper) sdk.CloudSDK {
					return cloudSdk
				},
				tokenProvider,
				func(transport http.RoundTripper, creds cloudpkg.Creds, tokenProvider cloudpkg.TokenProviderImpl, stack pkg.Stack) pkg.TokenProviderImpl {
					return stackTokenProvider
				},
				func(...formance.SDKOption) (sdk.StackSdkImpl, error) {
					return stacksdk, nil
				},
			)

			noopStep := resource.TestStep{
				Config: `
					provider "formancestack" {
						stack_id = "` + tc.stackId + `"
						organization_id = "` + tc.organizationId + `"
						uri = "` + tc.uri + `"
					}

					resource "formancestack_noop" "default" {}

				`,
			}

			if tc.expectedError != "" {
				noopStep.ExpectError = regexp.MustCompile(tc.expectedError)
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"formancestack": providerserver.NewProtocol6WithError(stackProvider()),
				},
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version0_15_0),
				},
				Steps: []resource.TestStep{
					noopStep,
				},
			})
		})
	}

}
