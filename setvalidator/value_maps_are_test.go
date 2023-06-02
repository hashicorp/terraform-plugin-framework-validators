// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package setvalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
)

func TestValueMapsAreValidatorValidateSet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Set
		elementValidators   []validator.Map
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.SetValueMust(
				types.MapType{ElemType: types.StringType},
				[]attr.Value{
					types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("first"),
							"Key2": types.StringValue("second"),
						},
					),
					types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("third"),
							"key2": types.StringValue("fourth"),
						},
					),
				},
			),
		},
		"Set unknown": {
			val: types.SetUnknown(
				types.StringType,
			),
			elementValidators: []validator.Map{
				mapvalidator.SizeAtLeast(1),
			},
		},
		"Set null": {
			val: types.SetNull(
				types.StringType,
			),
			elementValidators: []validator.Map{
				mapvalidator.SizeAtLeast(1),
			},
		},
		"Set elements invalid": {
			val: types.SetValueMust(
				types.MapType{ElemType: types.StringType},
				[]attr.Value{
					types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("first"),
							// Map ordering is random in Go, avoid multiple keys
						},
					),
					types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("third"),
							// Map ordering is random in Go, avoid multiple keys
						},
					),
				},
			),
			elementValidators: []validator.Map{
				mapvalidator.SizeAtLeast(3),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("first"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value({\"key1\":\"first\"})] map must contain at least 3 elements, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("third"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value({\"key1\":\"third\"})] map must contain at least 3 elements, got: 1",
				),
			},
		},
		"Set elements invalid for multiple validator": {
			val: types.SetValueMust(
				types.MapType{ElemType: types.StringType},
				[]attr.Value{
					types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("first"),
							// Map ordering is random in Go, avoid multiple keys
						},
					),
					types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("third"),
							// Map ordering is random in Go, avoid multiple keys
						},
					),
				},
			),
			elementValidators: []validator.Map{
				mapvalidator.SizeAtLeast(3),
				mapvalidator.SizeAtLeast(4),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("first"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value({\"key1\":\"first\"})] map must contain at least 3 elements, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("first"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value({\"key1\":\"first\"})] map must contain at least 4 elements, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("third"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value({\"key1\":\"third\"})] map must contain at least 3 elements, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("third"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value({\"key1\":\"third\"})] map must contain at least 4 elements, got: 1",
				),
			},
		},
		"Set elements valid": {
			val: types.SetValueMust(
				types.MapType{ElemType: types.StringType},
				[]attr.Value{
					types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("first"),
							"key2": types.StringValue("second"),
						},
					),
					types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("third"),
							"key2": types.StringValue("fourth"),
						},
					),
				},
			),
			elementValidators: []validator.Map{
				mapvalidator.SizeAtLeast(1),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.SetRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.val,
			}
			response := validator.SetResponse{}
			setvalidator.ValueMapsAre(testCase.elementValidators...).ValidateSet(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
