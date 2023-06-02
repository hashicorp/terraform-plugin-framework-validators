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

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
)

func TestValueInt64sAreValidatorValidateSet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Set
		elementValidators   []validator.Int64
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.SetValueMust(
				types.Int64Type,
				[]attr.Value{
					types.Int64Value(1),
					types.Int64Value(2),
				},
			),
		},
		"Set unknown": {
			val: types.SetUnknown(
				types.Int64Type,
			),
			elementValidators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"Set null": {
			val: types.SetNull(
				types.Int64Type,
			),
			elementValidators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"Set elements invalid": {
			val: types.SetValueMust(
				types.Int64Type,
				[]attr.Value{
					types.Int64Value(1),
					types.Int64Value(2),
				},
			),
			elementValidators: []validator.Int64{
				int64validator.AtLeast(3),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Int64Value(1)),
					"Invalid Attribute Value",
					"Attribute test[Value(1)] value must be at least 3, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Int64Value(2)),
					"Invalid Attribute Value",
					"Attribute test[Value(2)] value must be at least 3, got: 2",
				),
			},
		},
		"Set elements invalid for multiple validator": {
			val: types.SetValueMust(
				types.Int64Type,
				[]attr.Value{
					types.Int64Value(1),
					types.Int64Value(2),
				},
			),
			elementValidators: []validator.Int64{
				int64validator.AtLeast(3),
				int64validator.AtLeast(4),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Int64Value(1)),
					"Invalid Attribute Value",
					"Attribute test[Value(1)] value must be at least 3, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Int64Value(1)),
					"Invalid Attribute Value",
					"Attribute test[Value(1)] value must be at least 4, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Int64Value(2)),
					"Invalid Attribute Value",
					"Attribute test[Value(2)] value must be at least 3, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Int64Value(2)),
					"Invalid Attribute Value",
					"Attribute test[Value(2)] value must be at least 4, got: 2",
				),
			},
		},
		"Set elements valid": {
			val: types.SetValueMust(
				types.Int64Type,
				[]attr.Value{
					types.Int64Value(1),
					types.Int64Value(2),
				},
			),
			elementValidators: []validator.Int64{
				int64validator.AtLeast(1),
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
			setvalidator.ValueInt64sAre(testCase.elementValidators...).ValidateSet(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
