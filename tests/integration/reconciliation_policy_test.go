package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func TestReconciliationPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	cloudSdk := sdk.NewMockCloudSDK(ctrl)
	tokenProvider, _ := testprovider.NewMockTokenProvider(ctrl)
	stackTokenProvider := pkg.NewMockTokenProviderImpl(ctrl)
	stacksdk := sdk.NewMockStackSdkImpl(ctrl)
	reconciliationSdk := sdk.NewMockReconciliationSdkImpl(ctrl)
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
		func(...formance.SDKOption) sdk.StackSdkImpl {
			return stacksdk
		},
	)

	stacksdk.EXPECT().GetVersions(gomock.Any()).Return(&operations.GetVersionsResponse{
		GetVersionsResponse: &shared.GetVersionsResponse{
			Versions: []shared.Version{
				{
					Name:    "reconciliation",
					Version: "develop",
					Health:  true,
				},
			},
		},
	}, nil).AnyTimes()
	stacksdk.EXPECT().Reconciliation().Return(reconciliationSdk).AnyTimes()

	qry := `{
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

	m := make(map[string]any)
	require.NoError(t, json.Unmarshal([]byte(qry), &m), "Failed to unmarshal ledger query")
	policyId := uuid.NewString()
	policy := shared.Policy{
		Name:           "Test Policy",
		LedgerName:     "test-ledger",
		PaymentsPoolID: "test-payments-pool",
		ID:             policyId,
		LedgerQuery:    m,
		CreatedAt:      time.Now(),
	}
	// Init state
	reconciliationSdk.EXPECT().CreatePolicy(gomock.Any(), shared.PolicyRequest{
		LedgerName:     policy.LedgerName,
		Name:           policy.Name,
		PaymentsPoolID: policy.PaymentsPoolID,
		LedgerQuery:    policy.LedgerQuery,
	}).Return(&operations.CreatePolicyResponse{
		PolicyResponse: &shared.PolicyResponse{
			Data: policy,
		},
	}, nil)

	// refresh state deletion
	reconciliationSdk.EXPECT().GetPolicy(gomock.Any(), operations.GetPolicyRequest{
		PolicyID: policyId,
	}).Return(&operations.GetPolicyResponse{
		PolicyResponse: &shared.PolicyResponse{
			Data: policy,
		},
	}, nil)

	reconciliationSdk.EXPECT().DeletePolicy(gomock.Any(), operations.DeletePolicyRequest{
		PolicyID: policyId,
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

					resource "stack_reconciliation_policy" "policy" {
						ledger_name = "test-ledger"
						name = "Test Policy"
						payments_pool_id = "test-payments-pool"
						ledger_query = {
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
						}
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("stack_reconciliation_policy.policy", tfjsonpath.New("id"), knownvalue.StringExact(policyId)),
					statecheck.ExpectKnownValue("stack_reconciliation_policy.policy", tfjsonpath.New("ledger_name"), knownvalue.StringExact(policy.LedgerName)),
					statecheck.ExpectKnownValue("stack_reconciliation_policy.policy", tfjsonpath.New("name"), knownvalue.StringExact(policy.Name)),
					statecheck.ExpectKnownValue("stack_reconciliation_policy.policy", tfjsonpath.New("payments_pool_id"), knownvalue.StringExact(policy.PaymentsPoolID)),
					statecheck.ExpectKnownValue("stack_reconciliation_policy.policy", tfjsonpath.New("ledger_query"), knownvalue.MapExact(
						map[string]knownvalue.Check{
							"$and": knownvalue.TupleExact(
								[]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"$match": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"account": knownvalue.StringExact("accounts::pending"),
										}),
									}),
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"$or": knownvalue.TupleExact(
											[]knownvalue.Check{
												knownvalue.ObjectExact(map[string]knownvalue.Check{
													"$gte": knownvalue.ObjectExact(map[string]knownvalue.Check{
														"balance": knownvalue.Int64Exact(1000),
													}),
												}),
												knownvalue.ObjectExact(map[string]knownvalue.Check{
													"$lte": knownvalue.ObjectExact(map[string]knownvalue.Check{
														"balance": knownvalue.Int64Exact(500),
													}),
												}),
											},
										),
									}),
								}),
						}),
					),
				},
			},
		},
	})
}
