package resources_test

import (
	"fmt"
	"testing"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/pointer"
	"github.com/formancehq/terraform-provider-stack/internal/resources"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestPaymentsCreateConfigFromModel(t *testing.T) {
	t.Parallel()

	type testCase struct {
		plan              resources.PaymentsConnectorsModel
		expectedConnector operations.V3InstallConnectorRequest
	}

	for _, tc := range []testCase{
		{
			plan: resources.PaymentsConnectorsModel{
				Credentials: types.DynamicValue(types.MapValueMust(types.DynamicType, map[string]attr.Value{
					"apiKey": types.DynamicValue(types.StringValue("my-api-key")),
				})),
				Config: types.DynamicValue(types.MapValueMust(types.DynamicType, map[string]attr.Value{
					"endpoint":      types.DynamicValue(types.StringValue("https://api.example.com")),
					"name":          types.DynamicValue(types.StringValue("Example Connector")),
					"pageSize":      types.DynamicValue(types.Int64Value(100)),
					"pollingPeriod": types.DynamicValue(types.StringValue("5m")),
					"provider":      types.DynamicValue(types.StringValue("Generic")),
				})),
			},
			expectedConnector: operations.V3InstallConnectorRequest{
				V3InstallConnectorRequest: &shared.V3InstallConnectorRequest{
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
		},
		{
			plan: resources.PaymentsConnectorsModel{
				Credentials: types.DynamicValue(types.MapValueMust(types.DynamicType, map[string]attr.Value{
					"apiKey":          types.DynamicValue(types.StringValue("api-key-value")),
					"webhookPassword": types.DynamicValue(types.StringValue("webhook-password")),
				})),
				Config: types.DynamicValue(types.MapValueMust(types.DynamicType, map[string]attr.Value{
					"name":               types.DynamicValue(types.StringValue("Example Connector")),
					"pageSize":           types.DynamicValue(types.Int64Value(50)),
					"pollingPeriod":      types.DynamicValue(types.StringValue("2m")),
					"provider":           types.DynamicValue(types.StringValue("Adyen")),
					"companyID":          types.DynamicValue(types.StringValue("company-id-value")),
					"liveEndpointPrefix": types.DynamicValue(types.StringValue("https://live.example.com")),
					"webhookUsername":    types.DynamicValue(types.StringValue("webhook-username")),
				})),
			},
			expectedConnector: operations.V3InstallConnectorRequest{
				V3InstallConnectorRequest: &shared.V3InstallConnectorRequest{
					Type: "Adyen",
					V3AdyenConfig: &shared.V3AdyenConfig{
						APIKey:             "api-key-value",
						Name:               "Example Connector",
						PageSize:           pointer.For(int64(50)),
						PollingPeriod:      pointer.For("2m"),
						Provider:           pointer.For("Adyen"),
						CompanyID:          "company-id-value",
						LiveEndpointPrefix: pointer.For("https://live.example.com"),
						WebhookPassword:    pointer.For("webhook-password"),
						WebhookUsername:    pointer.For("webhook-username"),
					},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("%s %s", t.Name(), "1"), func(t *testing.T) {
			req, err := tc.plan.CreateConfig()
			if err != nil {
				t.Fatalf("failed to create config: %v", err)
			}

			diff := cmp.Diff(req, tc.expectedConnector)
			require.Empty(t, diff, "unexpected difference in connector config: %s", diff)
		})
	}
}

func TestExtractKeys(t *testing.T) {
	t.Parallel()
	type testCase struct {
		m            map[string]attr.Value
		expectedKeys []string
	}

	for _, tc := range []testCase{
		{
			m: map[string]attr.Value{
				"v1": types.DynamicNull(),
				"v2": types.DynamicNull(),
			},
			expectedKeys: []string{"v1", "v2"},
		},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			keys := resources.ExtractKeys(tc.m)

			require.ElementsMatch(t, tc.expectedKeys, keys)
		})
	}
}

func TestSanitizeUnknownKeys(t *testing.T) {
	t.Parallel()

	type testCase struct {
		expectedMap map[string]attr.Value
		m           map[string]attr.Value
		allowedKeys []string
	}

	for _, tc := range []testCase{
		{
			m: map[string]attr.Value{
				"v1": nil,
				"v2": nil,
				"v3": nil,
			},
			expectedMap: map[string]attr.Value{
				"v1": nil,
			},
			allowedKeys: []string{"v1"},
		},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			d := resources.SanitizeUnknownKeys(tc.m, tc.allowedKeys)
			require.Equal(t, tc.expectedMap, d)

		})
	}
}

func TestPaymentsStateFromRequest(t *testing.T) {
	t.Parallel()

	type testCase struct {
		request   *shared.V3GetConnectorConfigResponse
		fromState resources.PaymentsConnectorsModel
	}

	for _, tc := range []testCase{
		{
			request: &shared.V3GetConnectorConfigResponse{
				Data: shared.V3InstallConnectorRequest{
					V3AdyenConfig: &shared.V3AdyenConfig{
						APIKey:             "api-key-value",
						Name:               "Example Connector",
						PageSize:           pointer.For(int64(50)),
						PollingPeriod:      pointer.For("2m"),
						Provider:           pointer.For("Adyen"),
						CompanyID:          "company-id-value",
						LiveEndpointPrefix: pointer.For("https://live.example.com"),
						WebhookPassword:    pointer.For("webhook-password"),
						WebhookUsername:    pointer.For("webhook-username"),
					},
				},
			},
			fromState: resources.PaymentsConnectorsModel{
				ID: types.StringValue("somevalue"),
				Credentials: types.DynamicValue(types.ObjectValueMust(
					resources.GetMapTypeFromAttrTypes(map[string]attr.Value{
						"apiKey":          types.StringValue("api-key-value"),
						"webhookPassword": types.StringValue("webhook-password"),
					}), map[string]attr.Value{
						"apiKey":          types.StringValue("api-key-value"),
						"webhookPassword": types.StringValue("webhook-password"),
					})),
				Config: types.DynamicValue(types.ObjectValueMust(
					resources.GetMapTypeFromAttrTypes(map[string]attr.Value{
						"name":               types.StringValue("Example Connector"),
						"pageSize":           types.Int64Value(50),
						"pollingPeriod":      types.StringValue("2m"),
						"provider":           types.StringValue("Adyen"),
						"companyID":          types.StringValue("company-id-value"),
						"liveEndpointPrefix": types.StringValue("https://live.example.com"),
						"webhookUsername":    types.StringValue("webhook-username"),
					}), map[string]attr.Value{
						"name":               types.StringValue("Example Connector"),
						"pageSize":           types.Int64Value(50),
						"pollingPeriod":      types.StringValue("2m"),
						"provider":           types.StringValue("Adyen"),
						"companyID":          types.StringValue("company-id-value"),
						"liveEndpointPrefix": types.StringValue("https://live.example.com"),
						"webhookUsername":    types.StringValue("webhook-username"),
					})),
			},
		},
	} {
		t.Run(fmt.Sprintf("%s %s", t.Name(), *tc.request.Data.V3AdyenConfig.Provider), func(t *testing.T) {
			state, err := tc.fromState.StateFromRequest(&tc.request.Data)
			if err != nil {
				t.Fatalf("failed to create state from request: %v", err)
			}

			diff := cmp.Diff(state, tc.fromState)
			require.Empty(t, diff, "unexpected difference in state: %s", diff)
		})
	}
}
