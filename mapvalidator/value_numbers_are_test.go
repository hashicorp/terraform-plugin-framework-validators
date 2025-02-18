// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
)

func TestValueNumbersAreValidatorValidateMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Map
		elementValidators   []validator.Number
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.MapValueMust(
				types.NumberType,
				map[string]attr.Value{
					"key1": types.NumberValue(big.NewFloat(1.2)),
					"key2": types.NumberValue(big.NewFloat(2.4)),
				},
			),
		},
		"Map unknown": {
			val: types.MapUnknown(
				types.NumberType,
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(1.2)),
			},
		},
		"Map null": {
			val: types.MapNull(
				types.NumberType,
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(1.2)),
			},
		},
		"Map elements invalid": {
			val: types.MapValueMust(
				types.NumberType,
				map[string]attr.Value{
					"key1": types.NumberValue(big.NewFloat(1.2)),
					// Map ordering is random in Go, avoid multiple keys
				},
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(3.6)),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value Match",
					"Attribute test[\"key1\"] value must be one of: [\"3.6\"], got: 1.2",
				),
			},
		},
		"Map elements invalid for multiple validator": {
			val: types.MapValueMust(
				types.NumberType,
				map[string]attr.Value{
					"key1": types.NumberValue(big.NewFloat(1.2)),
					// Map ordering is random in Go, avoid multiple keys
				},
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(3.6)),
				numbervalidator.OneOf(big.NewFloat(4.8)),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value Match",
					"Attribute test[\"key1\"] value must be one of: [\"3.6\"], got: 1.2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value Match",
					"Attribute test[\"key1\"] value must be one of: [\"4.8\"], got: 1.2",
				),
			},
		},
		"Map elements valid": {
			val: types.MapValueMust(
				types.NumberType,
				map[string]attr.Value{
					"key1": types.NumberValue(big.NewFloat(1.2)),
					"key2": types.NumberValue(big.NewFloat(2.4)),
				},
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(1.2), big.NewFloat(2.4)),
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.MapRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.val,
			}
			response := validator.MapResponse{}
			mapvalidator.ValueNumbersAre(testCase.elementValidators...).ValidateMap(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
