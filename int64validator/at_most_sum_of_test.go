package int64validator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAtMostSumOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val                            attr.Value
		attributesToSumPathExpressions path.Expressions
		requestConfigRaw               map[string]tftypes.Value
		expectError                    bool
	}
	tests := map[string]testCase{
		"not an Int64": {
			val:         types.Bool{Value: true},
			expectError: true,
		},
		"unknown Int64": {
			val: types.Int64{Unknown: true},
		},
		"null Int64": {
			val: types.Int64{Null: true},
		},
		"valid integer as Int64 more than sum of attributes": {
			val: types.Int64{Value: 11},
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, 5),
				"two": tftypes.NewValue(tftypes.Number, 5),
			},
			expectError: true,
		},
		"valid integer as Int64 equal to sum of attributes": {
			val: types.Int64{Value: 10},
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, 5),
				"two": tftypes.NewValue(tftypes.Number, 5),
			},
		},
		"valid integer as Int64 less than sum of attributes": {
			val: types.Int64{Value: 7},
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, 4),
				"two": tftypes.NewValue(tftypes.Number, 4),
			},
		},
		"valid integer as Int64 less than sum of attributes, when one summed attribute is null": {
			val: types.Int64{Value: 8},
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, nil),
				"two": tftypes.NewValue(tftypes.Number, 9),
			},
		},
		"valid integer as Int64 does not return error when all attributes are null": {
			val: types.Int64{Null: true},
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, nil),
				"two": tftypes.NewValue(tftypes.Number, nil),
			},
		},
		"valid integer as Int64 returns error when all attributes to sum are null": {
			val: types.Int64{Value: 1},
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, nil),
				"two": tftypes.NewValue(tftypes.Number, nil),
			},
			expectError: true,
		},
		"valid integer as Int64 less than sum of attributes, when one summed attribute is unknown": {
			val: types.Int64{Value: 8},
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
				"two": tftypes.NewValue(tftypes.Number, 9),
			},
		},
		"valid integer as Int64 does not return error when all attributes are unknown": {
			val: types.Int64{Unknown: true},
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
				"two": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			},
		},
		"valid integer as Int64 does not return error when all attributes to sum are unknown": {
			val: types.Int64{Value: 1},
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
				"two": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			},
		},
		"error when attribute to sum is not Number": {
			val: types.Int64{Value: 9},
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Bool, true),
				"two": tftypes.NewValue(tftypes.Number, 9),
			},
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := tfsdk.ValidateAttributeRequest{
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
				AttributeConfig:         test.val,
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{}, test.requestConfigRaw),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {Type: types.Int64Type},
							"one":  {Type: types.Int64Type},
							"two":  {Type: types.Int64Type},
						},
					},
				},
			}

			response := tfsdk.ValidateAttributeResponse{}

			AtMostSumOf(test.attributesToSumPathExpressions...).Validate(context.Background(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
