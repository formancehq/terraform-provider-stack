package integration_test

import (
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

	"github.com/formancehq/go-libs/v3/logging"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/formancehq/terraform-provider-stack/internal/server"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"
)

func TestPaymentsPool(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	cloudSdk := sdk.NewMockCloudSDK(ctrl)
	tokenProvider, _ := cloudpkg.NewMockTokenProvider(ctrl)
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

	paymentsSdk.EXPECT().RemoveAccountFromPool(gomock.Any(), gomock.Cond(func(r operations.V3RemoveAccountFromPoolRequest) bool {
		return r.PoolID == poolId && r.AccountID == "account2"
	})).Return(nil, nil)

	// refresh state deletion
	firstPool.PoolAccounts = []string{"account1"}
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
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"formancestack": providerserver.NewProtocol6WithError(stackProvider()),
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
					provider "formancestack" {
						stack_id = "` + stackId + `"
						organization_id = "` + organizationId + `"
						uri = "` + fmt.Sprintf("https://%s-%s.formance.cloud/api", organizationId, stackId) + `"
					}

					resource "formancestack_payments_pool" "default" {
						name = "Example Pool"
						accounts_ids = [
							"account1",
							"account2",
						]
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"formancestack_payments_pool.default",
						tfjsonpath.New("id"),
						knownvalue.StringExact(poolId),
					),
					statecheck.ExpectKnownValue(
						"formancestack_payments_pool.default",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Example Pool"),
					),
					statecheck.ExpectKnownValue(
						"formancestack_payments_pool.default",
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
					provider "formancestack" {
						stack_id = "` + stackId + `"
						organization_id = "` + organizationId + `"
						uri = "` + fmt.Sprintf("https://%s-%s.formance.cloud/api", organizationId, stackId) + `"
					}

					resource "formancestack_payments_pool" "default" {
						name = "Example Pool"
						accounts_ids = [
							"account1",
						]
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"formancestack_payments_pool.default",
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
