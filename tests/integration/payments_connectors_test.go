package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/pkg/testprovider"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/formancehq/terraform-provider-stack/internal/server"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"
)

func TestPaymentsConnectors(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	cloudSdk := sdk.NewMockCloudSDK(ctrl)
	tokenProvider, _ := testprovider.NewMockTokenProvider(ctrl)
	stackTokenProvider := pkg.NewMockTokenProviderImpl(ctrl)
	stacksdk := sdk.NewMockStackSdkImpl(ctrl)
	paymentsSdk := sdk.NewMockPaymentsSdkImpl(ctrl)
	stackId := uuid.NewString()
	organizationId := uuid.NewString()

	stackProvider := server.NewStackProvider(
		logging.Testing().WithField("test", "payments_connectors"),
		server.FormanceStackEndpoint("dummy-endpoint"),
		server.FormanceStackClientId("organization_dummy-client-id"),
		server.FormanceStackClientSecret("dummy-client-secret"),
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

	stacksdk.EXPECT().GetVersions(gomock.Any()).Return(&operations.GetVersionsResponse{
		GetVersionsResponse: &shared.GetVersionsResponse{
			Versions: []shared.Version{
				{
					Name:    "payments",
					Version: "develop",
					Health:  true,
				},
			},
		},
	}, nil).AnyTimes()
	stacksdk.EXPECT().Payments().Return(paymentsSdk).AnyTimes()

	connectorId := uuid.NewString()
	// Init state
	paymentsSdk.EXPECT().CreateConnector(gomock.Any(), gomock.Cond(func(req operations.V3InstallConnectorRequest) bool {
		return true
	})).Return(&operations.V3InstallConnectorResponse{
		V3InstallConnectorResponse: &shared.V3InstallConnectorResponse{
			Data: connectorId,
		},
	}, nil)

	// Refresh state creation
	paymentsSdk.EXPECT().GetConnector(gomock.Any(), operations.V3GetConnectorConfigRequest{
		ConnectorID: connectorId,
	}).Return(&operations.V3GetConnectorConfigResponse{
		V3GetConnectorConfigResponse: &shared.V3GetConnectorConfigResponse{
			Data: shared.V3InstallConnectorRequest{
				Type: "Generic",
				V3GenericConfig: &shared.V3GenericConfig{
					APIKey:        "my-api-key",
					Endpoint:      "https://api.example.com",
					Name:          "Example Connector",
					PageSize:      pointer.For(int64(100)),
					PollingPeriod: pointer.For("5m"),
					Provider:      pointer.For("Generic"),
				},
			},
		},
	}, nil)

	// Refresh state update
	paymentsSdk.EXPECT().GetConnector(gomock.Any(), operations.V3GetConnectorConfigRequest{
		ConnectorID: connectorId,
	}).Return(&operations.V3GetConnectorConfigResponse{
		V3GetConnectorConfigResponse: &shared.V3GetConnectorConfigResponse{
			Data: shared.V3InstallConnectorRequest{
				Type: "Generic",
				V3GenericConfig: &shared.V3GenericConfig{
					APIKey:        "my-api-key",
					Endpoint:      "https://api.example.com",
					Name:          "Example Connector",
					PageSize:      pointer.For(int64(100)),
					PollingPeriod: pointer.For("5m"),
					Provider:      pointer.For("Generic"),
				},
			},
		},
	}, nil)

	paymentsSdk.EXPECT().UpdateConnector(gomock.Any(), gomock.Cond(func(r operations.V3UpdateConnectorConfigRequest) bool {
		return r.ConnectorID == connectorId
	})).Return(nil, nil)

	// refresh state deletion
	paymentsSdk.EXPECT().GetConnector(gomock.Any(), operations.V3GetConnectorConfigRequest{
		ConnectorID: connectorId,
	}).Return(&operations.V3GetConnectorConfigResponse{
		V3GetConnectorConfigResponse: &shared.V3GetConnectorConfigResponse{
			Data: shared.V3InstallConnectorRequest{
				Type: "Generic",
				V3GenericConfig: &shared.V3GenericConfig{
					APIKey:        "new-api-key",
					Endpoint:      "https://new-endpoint.com",
					Name:          "New Example Connector",
					PageSize:      pointer.For(int64(200)),
					PollingPeriod: pointer.For("10m"),
					Provider:      pointer.For("Generic"),
				},
			},
		},
	}, nil)

	paymentsSdk.EXPECT().DeleteConnector(gomock.Any(), operations.V3UninstallConnectorRequest{
		ConnectorID: connectorId,
	}).Return(nil, nil)

	// testCases
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stack": providerserver.NewProtocol6WithError(stackProvider()),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					provider "stack" {
						stack_id = "` + stackId + `"
						organization_id = "` + organizationId + `"
						uri = "` + fmt.Sprintf("https://%s-%s.formance.cloud/api", organizationId, stackId) + `"
					}

					resource "stack_payments_connectors" "generic" {
						credentials = {
							apiKey = "my-api-key"
						}

						config = {
							endpoint = "https://api.example.com"
							name = "Example Connector"
							pageSize = 100
							pollingPeriod = "5m"
							provider = "Generic"
						}
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("stack_payments_connectors.generic", tfjsonpath.New("credentials"), knownvalue.ObjectExact(
						map[string]knownvalue.Check{
							"apiKey": knownvalue.StringExact("my-api-key"),
						},
					)),
					statecheck.ExpectKnownValue("stack_payments_connectors.generic", tfjsonpath.New("config"), knownvalue.ObjectExact(
						map[string]knownvalue.Check{
							"endpoint":      knownvalue.StringExact("https://api.example.com"),
							"name":          knownvalue.StringExact("Example Connector"),
							"pageSize":      knownvalue.Int64Exact(100),
							"pollingPeriod": knownvalue.StringExact("5m"),
							"provider":      knownvalue.StringExact("Generic"),
						},
					)),
					statecheck.ExpectKnownValue("stack_payments_connectors.generic", tfjsonpath.New("id"), knownvalue.StringExact(connectorId)),
				},
			},
			{
				Config: `
					provider "stack" {
						stack_id = "` + stackId + `"
						organization_id = "` + organizationId + `"
						uri = "` + fmt.Sprintf("https://%s-%s.formance.cloud/api", organizationId, stackId) + `"
					}
					resource "stack_payments_connectors" "generic" {
						credentials = {
							apiKey = "new-api-key"
						}

						config = {
							endpoint = "https://new-endpoint.com"
							name = "New Example Connector"
							pageSize = 200
							pollingPeriod = "10m"
							provider = "Generic"
						}
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("stack_payments_connectors.generic", tfjsonpath.New("credentials"), knownvalue.ObjectExact(
						map[string]knownvalue.Check{
							"apiKey": knownvalue.StringExact("new-api-key"),
						},
					)),
					statecheck.ExpectKnownValue("stack_payments_connectors.generic", tfjsonpath.New("config"), knownvalue.ObjectExact(
						map[string]knownvalue.Check{
							"endpoint":      knownvalue.StringExact("https://new-endpoint.com"),
							"name":          knownvalue.StringExact("New Example Connector"),
							"pageSize":      knownvalue.Int64Exact(200),
							"pollingPeriod": knownvalue.StringExact("10m"),
							"provider":      knownvalue.StringExact("Generic"),
						},
					)),
					statecheck.ExpectKnownValue("stack_payments_connectors.generic", tfjsonpath.New("id"), knownvalue.StringExact(connectorId)),
				},
			},
		},
	})
}
