// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int32validator

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

func TestAtMostSumOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val                            types.Int32
		attributesToSumPathExpressions path.Expressions
		requestConfigRaw               map[string]tftypes.Value
		expectError                    bool
	}
	tests := map[string]testCase{
		"unknown Int32": {
			val: types.Int32Unknown(),
		},
		"null Int32": {
			val: types.Int32Null(),
		},
		"valid integer as Int32 more than sum of attributes": {
			val: types.Int32Value(11),
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
		"valid integer as Int32 equal to sum of attributes": {
			val: types.Int32Value(10),
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, 5),
				"two": tftypes.NewValue(tftypes.Number, 5),
			},
		},
		"valid integer as Int32 less than sum of attributes": {
			val: types.Int32Value(7),
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, 4),
				"two": tftypes.NewValue(tftypes.Number, 4),
			},
		},
		"valid integer as Int32 less than sum of attributes, when one summed attribute is null": {
			val: types.Int32Value(8),
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, nil),
				"two": tftypes.NewValue(tftypes.Number, 9),
			},
		},
		"valid integer as Int32 does not return error when all attributes are null": {
			val: types.Int32Null(),
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, nil),
				"two": tftypes.NewValue(tftypes.Number, nil),
			},
		},
		"valid integer as Int32 returns error when all attributes to sum are null": {
			val: types.Int32Value(1),
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
		"valid integer as Int32 less than sum of attributes, when one summed attribute is unknown": {
			val: types.Int32Value(8),
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
				"two": tftypes.NewValue(tftypes.Number, 9),
			},
		},
		"valid integer as Int32 does not return error when all attributes are unknown": {
			val: types.Int32Unknown(),
			attributesToSumPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
				"two": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			},
		},
		"valid integer as Int32 does not return error when all attributes to sum are unknown": {
			val: types.Int32Value(1),
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
			val: types.Int32Value(9),
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

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.Int32Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{}, test.requestConfigRaw),
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"test": schema.Int32Attribute{},
							"one":  schema.Int32Attribute{},
							"two":  schema.Int32Attribute{},
						},
					},
				},
			}

			response := validator.Int32Response{}

			AtMostSumOf(test.attributesToSumPathExpressions...).ValidateInt32(context.Background(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
