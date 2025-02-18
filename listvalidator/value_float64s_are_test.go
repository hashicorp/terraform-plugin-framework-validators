// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
)

func TestValueFloat64sAreValidatorValidateList(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.List
		elementValidators   []validator.Float64
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.ListValueMust(
				types.Float64Type,
				[]attr.Value{
					types.Float64Value(1),
					types.Float64Value(2),
				},
			),
		},
		"List unknown": {
			val: types.ListUnknown(
				types.Float64Type,
			),
			elementValidators: []validator.Float64{
				float64validator.AtLeast(1),
			},
		},
		"List null": {
			val: types.ListNull(
				types.Float64Type,
			),
			elementValidators: []validator.Float64{
				float64validator.AtLeast(1),
			},
		},
		"List elements invalid": {
			val: types.ListValueMust(
				types.Float64Type,
				[]attr.Value{
					types.Float64Value(1),
					types.Float64Value(2),
				},
			),
			elementValidators: []validator.Float64{
				float64validator.AtLeast(3),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value",
					"Attribute test[0] value must be at least 3.000000, got: 1.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(1),
					"Invalid Attribute Value",
					"Attribute test[1] value must be at least 3.000000, got: 2.000000",
				),
			},
		},
		"List elements invalid for multiple validator": {
			val: types.ListValueMust(
				types.Float64Type,
				[]attr.Value{
					types.Float64Value(1),
					types.Float64Value(2),
				},
			),
			elementValidators: []validator.Float64{
				float64validator.AtLeast(3),
				float64validator.AtLeast(4),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value",
					"Attribute test[0] value must be at least 3.000000, got: 1.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value",
					"Attribute test[0] value must be at least 4.000000, got: 1.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(1),
					"Invalid Attribute Value",
					"Attribute test[1] value must be at least 3.000000, got: 2.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(1),
					"Invalid Attribute Value",
					"Attribute test[1] value must be at least 4.000000, got: 2.000000",
				),
			},
		},
		"List elements valid": {
			val: types.ListValueMust(
				types.Float64Type,
				[]attr.Value{
					types.Float64Value(1),
					types.Float64Value(2),
				},
			),
			elementValidators: []validator.Float64{
				float64validator.AtLeast(1),
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.ListRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.val,
			}
			response := validator.ListResponse{}
			listvalidator.ValueFloat64sAre(testCase.elementValidators...).ValidateList(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
