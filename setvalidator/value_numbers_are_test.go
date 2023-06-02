// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package setvalidator_test

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

	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
)

func TestValueNumbersAreValidatorValidateSet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Set
		elementValidators   []validator.Number
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.SetValueMust(
				types.NumberType,
				[]attr.Value{
					types.NumberValue(big.NewFloat(1.2)),
					types.NumberValue(big.NewFloat(2.4)),
				},
			),
		},
		"Set unknown": {
			val: types.SetUnknown(
				types.NumberType,
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(1.2)),
			},
		},
		"Set null": {
			val: types.SetNull(
				types.NumberType,
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(1.2)),
			},
		},
		"Set elements invalid": {
			val: types.SetValueMust(
				types.NumberType,
				[]attr.Value{
					types.NumberValue(big.NewFloat(1.2)),
					types.NumberValue(big.NewFloat(2.4)),
				},
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(3.6)),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.NumberValue(big.NewFloat(1.2))),
					"Invalid Attribute Value Match",
					"Attribute test[Value(1.2)] value must be one of: [\"3.6\"], got: 1.2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.NumberValue(big.NewFloat(2.4))),
					"Invalid Attribute Value Match",
					"Attribute test[Value(2.4)] value must be one of: [\"3.6\"], got: 2.4",
				),
			},
		},
		"Set elements invalid for multiple validator": {
			val: types.SetValueMust(
				types.NumberType,
				[]attr.Value{
					types.NumberValue(big.NewFloat(1.2)),
					types.NumberValue(big.NewFloat(2.4)),
				},
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(3.6)),
				numbervalidator.OneOf(big.NewFloat(4.8)),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.NumberValue(big.NewFloat(1.2))),
					"Invalid Attribute Value Match",
					"Attribute test[Value(1.2)] value must be one of: [\"3.6\"], got: 1.2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.NumberValue(big.NewFloat(1.2))),
					"Invalid Attribute Value Match",
					"Attribute test[Value(1.2)] value must be one of: [\"4.8\"], got: 1.2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.NumberValue(big.NewFloat(2.4))),
					"Invalid Attribute Value Match",
					"Attribute test[Value(2.4)] value must be one of: [\"3.6\"], got: 2.4",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.NumberValue(big.NewFloat(2.4))),
					"Invalid Attribute Value Match",
					"Attribute test[Value(2.4)] value must be one of: [\"4.8\"], got: 2.4",
				),
			},
		},
		"Set elements valid": {
			val: types.SetValueMust(
				types.NumberType,
				[]attr.Value{
					types.NumberValue(big.NewFloat(1.2)),
					types.NumberValue(big.NewFloat(2.4)),
				},
			),
			elementValidators: []validator.Number{
				numbervalidator.OneOf(big.NewFloat(1.2), big.NewFloat(2.4)),
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
			setvalidator.ValueNumbersAre(testCase.elementValidators...).ValidateSet(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
