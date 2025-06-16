package resources

import (
	"fmt"
	"testing"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/v3/collectionutils"
	"github.com/formancehq/go-libs/v3/logging"
	"github.com/formancehq/terraform-provider-stack/internal"
	"github.com/formancehq/terraform-provider-stack/internal/server/sdk"
	"github.com/formancehq/terraform-provider-stack/pkg"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLedgerLedgerCreate(t *testing.T) {
	t.Parallel()

	r := NewLedger(logging.Testing())().(resource.ResourceWithConfigure)

	ctrl := gomock.NewController(t)
	stack := sdk.NewMockStackSdkImpl(ctrl)
	ledger := sdk.NewMockLedgerSdkImpl(ctrl)
	stack.EXPECT().Ledger().Return(ledger).AnyTimes()

	store := internal.Store{
		Stack:        pkg.Stack{},
		StackSdkImpl: stack,
	}

	req := resource.ConfigureRequest{
		ProviderData: store,
	}
	res := resource.ConfigureResponse{
		Diagnostics: diag.Diagnostics{},
	}

	stack.EXPECT().GetVersions(gomock.Any()).Return(&operations.GetVersionsResponse{
		GetVersionsResponse: &shared.GetVersionsResponse{
			Versions: []shared.Version{
				{
					Health:  true,
					Name:    "ledger",
					Version: "develop",
				},
			},
		},
	}, nil)

	r.Configure(logging.TestingContext(), req, &res)

	require.Len(t, res.Diagnostics, 0, "Expected no diagnostics after configuring the resource")

	values := map[string]tftypes.Value{
		"name":   tftypes.NewValue(tftypes.String, uuid.NewString()),
		"bucket": tftypes.NewValue(tftypes.String, uuid.NewString()),
		"metadata": tftypes.NewValue(tftypes.Map{
			ElementType: tftypes.String,
		}, map[string]tftypes.Value{
			"key1": tftypes.NewValue(tftypes.String, "value1"),
			"key2": tftypes.NewValue(tftypes.String, "value2"),
		}),
	}
	createReq := resource.CreateRequest{
		Plan: tfsdk.Plan{
			Raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"name":   tftypes.String,
					"bucket": tftypes.String,
					"metadata": tftypes.Map{
						ElementType: tftypes.String,
					},
				},
			}, values),
			Schema: SchemaLedger,
		},
	}
	createRes := resource.CreateResponse{
		Diagnostics: diag.Diagnostics{},
		State: tfsdk.State{
			Schema: SchemaLedger,
		},
	}

	ledger.EXPECT().CreateLedger(gomock.Any(), gomock.Cond(func(req operations.V2CreateLedgerRequest) bool {
		var (
			name     string
			bucket   *string
			metadata map[string]tftypes.Value
		)

		require.NoError(t, values["name"].As(&name), "Expected name to be a string")
		require.NoError(t, values["bucket"].As(&bucket), "Expected bucket to be a string")
		require.NoError(t, values["metadata"].As(&metadata), "Expected metadata to be a map[string]Value")

		v := operations.V2CreateLedgerRequest{
			Ledger: name,
			V2CreateLedgerRequest: shared.V2CreateLedgerRequest{
				Bucket: bucket,
				Metadata: collectionutils.ConvertMap(metadata, func(v tftypes.Value) string {
					var str string

					require.NoError(t, v.As(&str), "Expected metadata value to be a string")
					return str
				}),
			},
		}

		diff := cmp.Diff(req, v)
		if diff != "" {
			fmt.Println(diff)
		}
		return diff == ""
	})).Return(&operations.V2CreateLedgerResponse{
		RawResponse: nil,
	}, nil)
	r.Create(logging.TestingContext(), createReq, &createRes)

	require.Len(t, createRes.Diagnostics, 0, "Expected no diagnostics after creating the resource")
}
