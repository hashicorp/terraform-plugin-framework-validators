// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package int64validator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAtLeastSumOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val                        types.Int64
		attributesToSumExpressions path.Expressions
		requestConfigRaw           map[string]tftypes.Value
		expectError                bool
	}
	tests := map[string]testCase{
		"unknown Int64": {
			val: types.Int64Unknown(),
		},
		"null Int64": {
			val: types.Int64Null(),
		},
		"valid integer as Int64 less than sum of attributes": {
			val: types.Int64Value(10),
			attributesToSumExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, 15),
				"two": tftypes.NewValue(tftypes.Number, 15),
			},
			expectError: true,
		},
		"valid integer as Int64 equal to sum of attributes": {
			val: types.Int64Value(10),
			attributesToSumExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, 5),
				"two": tftypes.NewValue(tftypes.Number, 5),
			},
		},
		"valid integer as Int64 greater than sum of attributes": {
			val: types.Int64Value(10),
			attributesToSumExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, 4),
				"two": tftypes.NewValue(tftypes.Number, 4),
			},
		},
		"valid integer as Int64 greater than sum of attributes, when one summed attribute is null": {
			val: types.Int64Value(10),
			attributesToSumExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, nil),
				"two": tftypes.NewValue(tftypes.Number, 9),
			},
		},
		"valid integer as Int64 does not return error when all attributes are null": {
			val: types.Int64Null(),
			attributesToSumExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, nil),
				"two": tftypes.NewValue(tftypes.Number, nil),
			},
		},
		"valid integer as Int64 returns error when all attributes to sum are null": {
			val: types.Int64Value(-1),
			attributesToSumExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, nil),
				"two": tftypes.NewValue(tftypes.Number, nil),
			},
			expectError: true,
		},
		"valid integer as Int64 greater than sum of attributes, when one summed attribute is unknown": {
			val: types.Int64Value(10),
			attributesToSumExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
				"two": tftypes.NewValue(tftypes.Number, 9),
			},
		},
		"valid integer as Int64 does not return error when all attributes are unknown": {
			val: types.Int64Unknown(),
			attributesToSumExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
				"two": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			},
		},
		"valid integer as Int64 does not return error when all attributes to sum are unknown": {
			val: types.Int64Value(-1),
			attributesToSumExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
				"two": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			},
		},
		"error when attribute to sum is not Number": {
			val: types.Int64Value(9),
			attributesToSumExpressions: path.Expressions{
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

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.Int64Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{}, test.requestConfigRaw),
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"test": schema.Int64Attribute{},
							"one":  schema.Int64Attribute{},
							"two":  schema.Int64Attribute{},
						},
					},
				},
			}

			response := validator.Int64Response{}

			AtLeastSumOf(test.attributesToSumExpressions...).ValidateInt64(context.Background(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
