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
	"go.opentelemetry.io/otel"

	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/go-libs/v3/pointer"
	cloudpkg "github.com/formancehq/terraform-provider-cloud/pkg"
	"github.com/formancehq/terraform-provider-cloud/pkg/testprovider"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/formancehq/terraform-provider-stack/internal/server"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"
)

func TestLedgerSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	cloudSdk := sdk.NewMockCloudSDK(ctrl)
	tokenProvider, _ := testprovider.NewMockTokenProvider(ctrl)
	stackTokenProvider := pkg.NewMockTokenProviderImpl(ctrl)
	stacksdk := sdk.NewMockStackSdkImpl(ctrl)
	ledgerSchema := sdk.NewMockLedgerSdkImpl(ctrl)
	stackId := uuid.NewString()
	organizationId := uuid.NewString()

	stackProvider := server.NewStackProvider(
		otel.GetTracerProvider(),

		logging.Testing().WithField("test", t.Name()),
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
					Name:    "ledger",
					Version: "develop",
					Health:  true,
				},
			},
		},
	}, nil).AnyTimes()
	stacksdk.EXPECT().Ledger().Return(ledgerSchema).AnyTimes()

	schema := map[string]shared.V2ChartSegment{
		"segment1": {
			DotSelf: &shared.DotSelf{},
		},
		"segment2": {
			DotMetadata: map[string]shared.V2ChartAccountMetadata{
				"test": {
					Default: pointer.For("test"),
				},
			},
		},
	}
	ledgerSchema.EXPECT().InsertSchema(gomock.Any(), gomock.Cond(func(op operations.V2InsertSchemaRequest) bool {
		return cmp.Diff(op, operations.V2InsertSchemaRequest{
			Ledger:  "test-ledger",
			Version: "v1.0.0",
			V2SchemaData: shared.V2SchemaData{
				Chart: schema,
			},
		}) != ""
	})).Return(&operations.V2InsertSchemaResponse{
		StatusCode: http.StatusOK,
	}, nil).Times(1)

	ledgerSchema.EXPECT().GetSchema(gomock.Any(), operations.V2GetSchemaRequest{
		Ledger:  "test-ledger",
		Version: "v1.0.0",
	}).Return(&operations.V2GetSchemaResponse{
		StatusCode: http.StatusOK,
		V2SchemaResponse: &shared.V2SchemaResponse{
			Data: shared.V2Schema{
				Version: "v1.0.0",
				Chart:   schema,
			},
		},
	}, nil).Times(2)

	schemaUpdated := map[string]shared.V2ChartSegment{
		"segment3": {
			DotSelf: &shared.DotSelf{},
		},
		"segment2": {
			DotMetadata: map[string]shared.V2ChartAccountMetadata{
				"test": {
					Default: pointer.For("test"),
				},
			},
		},
	}
	ledgerSchema.EXPECT().InsertSchema(gomock.Any(), gomock.Cond(func(op operations.V2InsertSchemaRequest) bool {
		return cmp.Diff(op, operations.V2InsertSchemaRequest{
			Ledger:  "test-ledger",
			Version: "v1.0.1",
			V2SchemaData: shared.V2SchemaData{
				Chart: schemaUpdated,
			},
		}) != ""
	})).Return(&operations.V2InsertSchemaResponse{
		StatusCode: http.StatusOK,
	}, nil)

	ledgerSchema.EXPECT().GetSchema(gomock.Any(), operations.V2GetSchemaRequest{
		Ledger:  "test-ledger",
		Version: "v1.0.1",
	}).Return(&operations.V2GetSchemaResponse{
		StatusCode: http.StatusOK,
		V2SchemaResponse: &shared.V2SchemaResponse{
			Data: shared.V2Schema{
				Version: "v1.0.1",
				Chart:   schemaUpdated,
			},
		},
	}, nil)
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

					resource "stack_ledger_schema" "default" {
						ledger = "test-ledger"
						version = "v1.0.0"
						schema = {
							"segment1": {
								".self": {}
							},
							"segment2": {
								".metadata": {
									"test": {
										"default": "test"
									}
								}
							}
						}
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"stack_ledger_schema.default",
						tfjsonpath.New(`schema`),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"segment1": knownvalue.ObjectExact(map[string]knownvalue.Check{
								".self": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							}),
							"segment2": knownvalue.ObjectExact(map[string]knownvalue.Check{
								".metadata": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"test": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"default": knownvalue.StringExact("test"),
									}),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"stack_ledger_schema.default",
						tfjsonpath.New("version"),
						knownvalue.StringExact("v1.0.0"),
					),
					statecheck.ExpectKnownValue(
						"stack_ledger_schema.default",
						tfjsonpath.New("ledger"),
						knownvalue.StringExact("test-ledger"),
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

					resource "stack_ledger_schema" "default" {
						ledger = "test-ledger"
						version = "v1.0.1"
						schema = {
							"segment3": {
								".self": {}
							},
							"segment2": {
								".metadata": {
									"test": {
										"default": "test"
									}
								}
							}
						}
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"stack_ledger_schema.default",
						tfjsonpath.New(`schema`),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"segment3": knownvalue.ObjectExact(map[string]knownvalue.Check{
								".self": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							}),
							"segment2": knownvalue.ObjectExact(map[string]knownvalue.Check{
								".metadata": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"test": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"default": knownvalue.StringExact("test"),
									}),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"stack_ledger_schema.default",
						tfjsonpath.New("version"),
						knownvalue.StringExact("v1.0.1"),
					),
					statecheck.ExpectKnownValue(
						"stack_ledger_schema.default",
						tfjsonpath.New("ledger"),
						knownvalue.StringExact("test-ledger"),
					),
				},
			},
		},
	})
}
