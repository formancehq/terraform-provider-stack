package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

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
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"

	"github.com/formancehq/go-libs/v3/logging"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/pkg/testprovider"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/formancehq/terraform-provider-stack/internal/server"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"
)

func TestPaymentsPool(t *testing.T) {
	ctrl := gomock.NewController(t)
	cloudSdk := sdk.NewMockCloudSDK(ctrl)
	tokenProvider, _ := testprovider.NewMockTokenProvider(ctrl)
	stackTokenProvider := pkg.NewMockTokenProviderImpl(ctrl)
	stacksdk := sdk.NewMockStackSdkImpl(ctrl)
	paymentsSdk := sdk.NewMockPaymentsSdkImpl(ctrl)
	stackId := uuid.NewString()
	organizationId := uuid.NewString()

	stackProvider := server.NewStackProvider(
		otel.GetTracerProvider(),

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

	query := `{
		"$and": [
			{
				"$match": {
					"account": "accounts::pending"
				}
			},
			{
				"$or": [
					{
						"$gte": {
							"balance": 1000
						}
					},
					{
						"$lte": {
							"balance": 500
						}
					}
				]
			}
		]
	}`
	queryAsMap := make(map[string]any)
	require.NoError(t, json.Unmarshal([]byte(query), &queryAsMap))

	queryUpdated := `{
		"$and": [
			{
				"$match": {
					"account": "accounts::another"
				}
			},
			{
				"$or": [
					{
						"$gte": {
							"balance": 1000
						}
					},
					{
						"$lte": {
							"balance": 500
						}
					}
				]
			}
		]
	}`
	queryUpdatedAsMap := make(map[string]any)
	require.NoError(t, json.Unmarshal([]byte(queryUpdated), &queryUpdatedAsMap))

	// Module and sdk expectations
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

	poolId := uuid.NewString()
	firstPool := shared.V3Pool{
		ID:           poolId,
		Name:         "Example Pool",
		PoolAccounts: []string{"account1", "account2"},
		Query:        queryAsMap,
		CreatedAt:    time.Now(),
	}

	// Init state
	paymentsSdk.EXPECT().CreatePool(gomock.Any(), gomock.Cond(func(req *shared.V3CreatePoolRequest) bool {
		return true
	})).Return(&operations.V3CreatePoolResponse{
		V3CreatePoolResponse: &shared.V3CreatePoolResponse{
			Data: poolId,
		},
	}, nil)

	// Refresh state creation
	paymentsSdk.EXPECT().GetPool(gomock.Any(), operations.V3GetPoolRequest{
		PoolID: poolId,
	}).Return(&operations.V3GetPoolResponse{
		V3GetPoolResponse: &shared.V3GetPoolResponse{
			Data: firstPool,
		},
	}, nil)

	// Refresh state update
	paymentsSdk.EXPECT().GetPool(gomock.Any(), operations.V3GetPoolRequest{
		PoolID: poolId,
	}).Return(&operations.V3GetPoolResponse{
		V3GetPoolResponse: &shared.V3GetPoolResponse{
			Data: firstPool,
		},
	}, nil)

	paymentsSdk.EXPECT().UpdatePool(gomock.Any(), gomock.Cond(func(op operations.V3UpdatePoolQueryRequest) bool {
		return strings.ReplaceAll(op.PoolID, "\"", "") == poolId && fmt.Sprintf("%v", op.V3UpdatePoolQueryRequest.Query) == fmt.Sprintf("%v", queryUpdatedAsMap)
	})).Return(&operations.V3UpdatePoolQueryResponse{
		StatusCode: 200,
	}, nil)

	paymentsSdk.EXPECT().RemoveAccountFromPool(gomock.Any(), gomock.Cond(func(r operations.V3RemoveAccountFromPoolRequest) bool {
		return r.PoolID == poolId && r.AccountID == "account2"
	})).Return(nil, nil)

	// refresh state deletion
	firstPool.PoolAccounts = []string{"account1"}
	firstPool.Query = queryUpdatedAsMap
	paymentsSdk.EXPECT().GetPool(gomock.Any(), operations.V3GetPoolRequest{
		PoolID: poolId,
	}).Return(&operations.V3GetPoolResponse{
		V3GetPoolResponse: &shared.V3GetPoolResponse{
			Data: firstPool,
		},
	}, nil)

	paymentsSdk.EXPECT().DeletePool(gomock.Any(), operations.V3DeletePoolRequest{
		PoolID: poolId,
	}).Return(nil, nil)

	// testCases
	resource.ParallelTest(t, resource.TestCase{
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

					resource "stack_payments_pool" "default" {
						name = "Example Pool"
						accounts_ids = [
							"account1",
							"account2",
						]
						query = ` + query + `
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"stack_payments_pool.default",
						tfjsonpath.New("id"),
						knownvalue.StringExact(poolId),
					),
					statecheck.ExpectKnownValue(
						"stack_payments_pool.default",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Example Pool"),
					),
					statecheck.ExpectKnownValue(
						"stack_payments_pool.default",
						tfjsonpath.New("accounts_ids"),
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.StringExact("account1"),
								knownvalue.StringExact("account2"),
							},
						),
					),
				},
			},
			{
				Config: `
					provider "stack" {
						stack_id = "` + stackId + `"
						organization_id = "` + organizationId + `"
						uri = "` + fmt.Sprintf("https://%s-%s.formance.cloud/api", organizationId, stackId) + `"
					}

					resource "stack_payments_pool" "default" {
						name = "Example Pool"
						accounts_ids = [
							"account1",
						]
						query = ` + queryUpdated + `
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"stack_payments_pool.default",
						tfjsonpath.New("accounts_ids"),
						knownvalue.ListExact(
							[]knownvalue.Check{
								knownvalue.StringExact("account1"),
							},
						),
					),
				},
			},
		},
	})
}
