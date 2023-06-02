// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator_test

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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestValueStringsAreValidatorValidateMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Map
		elementValidators   []validator.String
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("first"),
					"key2": types.StringValue("second"),
				},
			),
		},
		"Map unknown": {
			val: types.MapUnknown(
				types.StringType,
			),
			elementValidators: []validator.String{
				stringvalidator.LengthAtLeast(6),
			},
		},
		"Map null": {
			val: types.MapNull(
				types.StringType,
			),
			elementValidators: []validator.String{
				stringvalidator.LengthAtLeast(6),
			},
		},
		"Map elements invalid": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("first"),
					// Map ordering is random in Go, avoid multiple keys
				},
			),
			elementValidators: []validator.String{
				stringvalidator.LengthAtLeast(7),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value Length",
					"Attribute test[\"key1\"] string length must be at least 7, got: 5",
				),
			},
		},
		"Map elements invalid for multiple validator": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("first"),
					// Map ordering is random in Go, avoid multiple keys
				},
			),
			elementValidators: []validator.String{
				stringvalidator.LengthAtLeast(7),
				stringvalidator.LengthAtLeast(8),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value Length",
					"Attribute test[\"key1\"] string length must be at least 7, got: 5",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value Length",
					"Attribute test[\"key1\"] string length must be at least 8, got: 5",
				),
			},
		},
		"Map elements valid": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("first"),
					"key2": types.StringValue("second"),
				},
			),
			elementValidators: []validator.String{
				stringvalidator.LengthAtLeast(5),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.MapRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.val,
			}
			response := validator.MapResponse{}
			mapvalidator.ValueStringsAre(testCase.elementValidators...).ValidateMap(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
