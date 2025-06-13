package server_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/formancehq/go-libs/v3/logging"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	cloudsdk "github.com/formancehq/terraform-provider-cloud/sdk"
	"github.com/formancehq/terraform-provider-stack/internal/server"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"
	"go.uber.org/mock/gomock"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestProviderMetadata(t *testing.T) {
	t.Parallel()
	p := server.NewStackProvider(logging.Testing(),
		"https://app.formance.cloud/api",
		"client_id",
		"client_secret",
		http.DefaultTransport,
		sdk.NewCloudSDK(),
		cloudpkg.NewTokenProvider,
		func(transport http.RoundTripper, creds cloudpkg.Creds, tokenProvider cloudpkg.TokenProviderImpl, stack pkg.Stack) pkg.TokenProviderImpl {
			return pkg.NewTokenProvider(transport, creds, tokenProvider, stack)
		},
		sdk.NewStackSdk(),
	)()

	res := provider.MetadataResponse{}
	p.Metadata(logging.TestingContext(), provider.MetadataRequest{}, &res)

	require.Equal(t, res.TypeName, "formancestack")
	require.Equal(t, res.Version, "develop")
}

func TestProviderConfigure(t *testing.T) {
	t.Parallel()
	type testCase struct {
		ClientId     string
		ClientSecret string
		Endpoint     string
	}

	for _, tc := range []testCase{
		{
			ClientSecret: uuid.NewString(),
		},
		{
			ClientId: fmt.Sprintf("organization_%s", uuid.NewString()),
		},
		{
			Endpoint: uuid.NewString(),
		},
		{},
	} {
		t.Run(fmt.Sprintf("%s clientId %s clientSecret %s endpoint %s", t.Name(), tc.ClientId, tc.ClientSecret, tc.Endpoint), func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			cloudSdk := sdk.NewMockCloudSDK(ctrl)
			tokenProvider, mock := cloudpkg.NewMockTokenProvider(ctrl)
			stackTokenProvider := pkg.NewMockTokenProviderImpl(ctrl)
			stacksdk := sdk.NewMockStackSdkImpl(ctrl)

			p := server.NewStackProvider(
				logging.Testing(),
				"https://app.formance.cloud/api",
				"organization_client_id",
				"client_secret",
				http.DefaultTransport,
				func(creds cloudpkg.Creds, transport http.RoundTripper) sdk.CloudSDK {
					return cloudSdk
				},
				tokenProvider,
				func(transport http.RoundTripper, creds cloudpkg.Creds, tokenProvider cloudpkg.TokenProviderImpl, stack pkg.Stack) pkg.TokenProviderImpl {
					return stackTokenProvider
				},
				func(url, version string, transport http.RoundTripper, tp pkg.TokenProviderImpl) (sdk.StackSdkImpl, error) {
					return stacksdk, nil
				},
			)()
			stackId := uuid.NewString()
			organizationId := uuid.NewString()
			stackUri := fmt.Sprintf("https://%s-%s.formance.cloud/api", organizationId, stackId)

			clientId := tc.ClientId
			if clientId == "" {
				clientId = "organization_client_id"
			}
			if server.IsOrganizationClient(clientId) {
				cloudSdk.EXPECT().GetStack(gomock.Any(), organizationId, stackId).Return(&cloudsdk.CreateStackResponse{
					Data: &cloudsdk.Stack{
						State:  "ACTIVE",
						Status: "READY",
					},
				}, nil, nil)
			}

			res := provider.ConfigureResponse{
				Diagnostics: []diag.Diagnostic{},
			}

			cloudObj := tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"client_id":     tftypes.String,
					"client_secret": tftypes.String,
					"endpoint":      tftypes.String,
				},
			}
			p.Configure(logging.TestingContext(), provider.ConfigureRequest{
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"cloud":           cloudObj,
							"stack_id":        tftypes.String,
							"organization_id": tftypes.String,
							"uri":             tftypes.String,
							"expected_modules": tftypes.List{
								ElementType: tftypes.String,
							},
						},
					}, map[string]tftypes.Value{
						"stack_id":         tftypes.NewValue(tftypes.String, stackId),
						"organization_id":  tftypes.NewValue(tftypes.String, organizationId),
						"uri":              tftypes.NewValue(tftypes.String, stackUri),
						"expected_modules": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{}),
						"cloud": tftypes.NewValue(cloudObj, map[string]tftypes.Value{
							"client_id":     tftypes.NewValue(tftypes.String, tc.ClientId),
							"client_secret": tftypes.NewValue(tftypes.String, tc.ClientSecret),
							"endpoint":      tftypes.NewValue(tftypes.String, tc.Endpoint),
						}),
					}),
					Schema: server.SchemaStack,
				},
			}, &res)

			if !server.IsOrganizationClient(clientId) {
				require.NotEmpty(t, res.Diagnostics)
				return
			}

			require.Empty(t, res.Diagnostics)
			require.NotNil(t, res.ResourceData)
			require.IsType(t, stacksdk, res.ResourceData)
			require.NotNil(t, res.DataSourceData)
			require.IsType(t, stacksdk, res.DataSourceData)

			if tc.ClientId == "" {
				require.Equal(t, mock.ClientId(), "organization_client_id")
			} else {
				require.Equal(t, mock.ClientId(), tc.ClientId)
			}

			if tc.ClientSecret == "" {
				require.Equal(t, mock.ClientSecret(), "client_secret")
			} else {
				require.Equal(t, mock.ClientSecret(), tc.ClientSecret)
			}

			if tc.Endpoint == "" {
				require.Equal(t, mock.Endpoint(), "https://app.formance.cloud/api")
			} else {
				require.Equal(t, mock.Endpoint(), tc.Endpoint)
			}
		})
	}

}
