package resources_test

import (
	"testing"

	"github.com/formancehq/terraform-provider-stack/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestCreateConfig(t *testing.T) {
	t.Parallel()

	type testCase struct {
		ledgerQuery resources.DynamicObjectValue
		expectedErr error
	}

	for _, tc := range []testCase{
		{
			ledgerQuery: resources.NewDynamicObjectValue(nil),
			expectedErr: resources.ErrParseLedgerQuery,
		},
		{
			ledgerQuery: resources.NewDynamicObjectValue(map[string]attr.Value{}),
		},
		{
			ledgerQuery: resources.NewDynamicObjectValue(map[string]attr.Value{
				"toto": types.StringNull(),
			}),
			expectedErr: resources.ErrParseLedgerQuery,
		},
		{
			ledgerQuery: resources.NewDynamicObjectValue(map[string]attr.Value{
				"$and": types.StringNull(),
			}),
			expectedErr: resources.ErrParseLedgerQuery,
		},
		{
			ledgerQuery: resources.NewDynamicObjectValue(map[string]attr.Value{
				"$not": resources.NewDynamicObjectValue(map[string]attr.Value{
					"$and": resources.NewDynamicTupleValue(
						[]attr.Value{
							resources.NewDynamicObjectValue(
								map[string]attr.Value{
									"$match": resources.NewDynamicObjectValue(map[string]attr.Value{
										"account": types.StringValue("accounts::pending"),
									}).Value(),
								},
							).Value(),
							resources.NewDynamicObjectValue(
								map[string]attr.Value{
									"$or": resources.NewDynamicTupleValue(
										[]attr.Value{
											resources.NewDynamicObjectValue(
												map[string]attr.Value{
													"$gte": resources.NewDynamicObjectValue(
														map[string]attr.Value{
															"balance": types.Int64Value(1000),
														},
													).Value(),
												},
											).Value(),
											resources.NewDynamicObjectValue(
												map[string]attr.Value{
													"$match": resources.NewDynamicObjectValue(
														map[string]attr.Value{
															"metadata[category]": types.StringValue("gold"),
														},
													).Value(),
												},
											).Value(),
										},
									).Value(),
								},
							).Value(),
						},
					).Value(),
				}).Value(),
			}),
		},
	} {
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			model := &resources.ReconciliationPolicyModel{
				LedgerQuery: types.DynamicValue(tc.ledgerQuery.Value()),
			}

			_, err := model.CreateConfig()
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
