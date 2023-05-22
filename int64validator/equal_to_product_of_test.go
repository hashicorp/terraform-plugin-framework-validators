// Copyright (c) HashiCorp, Inc.
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

func TestEqualToProductOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val                                 types.Int64
		attributesToMultiplyPathExpressions path.Expressions
		requestConfigRaw                    map[string]tftypes.Value
		expectError                         bool
	}
	tests := map[string]testCase{
		"unknown Int64": {
			val: types.Int64Unknown(),
		},
		"null Int64": {
			val: types.Int64Null(),
		},
		"valid integer as Int64 more than product of attributes": {
			val: types.Int64Value(26),
			attributesToMultiplyPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, 5),
				"two": tftypes.NewValue(tftypes.Number, 5),
			},
			expectError: true,
		},
		"valid integer as Int64 less than product of attributes": {
			val: types.Int64Value(24),
			attributesToMultiplyPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, 5),
				"two": tftypes.NewValue(tftypes.Number, 5),
			},
			expectError: true,
		},
		"valid integer as Int64 equal to product of attributes": {
			val: types.Int64Value(25),
			attributesToMultiplyPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, 5),
				"two": tftypes.NewValue(tftypes.Number, 5),
			},
		},
		"validation skipped when one attribute is null": {
			val: types.Int64Value(10),
			attributesToMultiplyPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, nil),
				"two": tftypes.NewValue(tftypes.Number, 8),
			},
		},
		"validation skipped when all attributes are null": {
			val: types.Int64Null(),
			attributesToMultiplyPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, nil),
				"two": tftypes.NewValue(tftypes.Number, nil),
			},
		},
		"validation skipped when all attributes to multiply are null": {
			val: types.Int64Value(1),
			attributesToMultiplyPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, nil),
				"two": tftypes.NewValue(tftypes.Number, nil),
			},
		},
		"validation skipped when one attribute is unknown": {
			val: types.Int64Value(10),
			attributesToMultiplyPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
				"two": tftypes.NewValue(tftypes.Number, 8),
			},
		},
		"validation skipped when all attributes are unknown": {
			val: types.Int64Unknown(),
			attributesToMultiplyPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
				"two": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			},
		},
		"validation skipped when all attributes to multiply are unknown": {
			val: types.Int64Value(1),
			attributesToMultiplyPathExpressions: path.Expressions{
				path.MatchRoot("one"),
				path.MatchRoot("two"),
			},
			requestConfigRaw: map[string]tftypes.Value{
				"one": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
				"two": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
			},
		},
		"error when attribute to multiply is not Number": {
			val: types.Int64Value(9),
			attributesToMultiplyPathExpressions: path.Expressions{
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

			EqualToProductOf(test.attributesToMultiplyPathExpressions...).ValidateInt64(context.Background(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
