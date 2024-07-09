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

	"github.com/hashicorp/terraform-plugin-framework-validators/float32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
)

func TestValueFloat32sAreValidatorValidateSet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Set
		elementValidators   []validator.Float32
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.SetValueMust(
				types.Float32Type,
				[]attr.Value{
					types.Float32Value(1),
					types.Float32Value(2),
				},
			),
		},
		"Set unknown": {
			val: types.SetUnknown(
				types.Float32Type,
			),
			elementValidators: []validator.Float32{
				float32validator.AtLeast(1),
			},
		},
		"Set null": {
			val: types.SetNull(
				types.Float32Type,
			),
			elementValidators: []validator.Float32{
				float32validator.AtLeast(1),
			},
		},
		"Set elements invalid": {
			val: types.SetValueMust(
				types.Float32Type,
				[]attr.Value{
					types.Float32Value(1),
					types.Float32Value(2),
				},
			),
			elementValidators: []validator.Float32{
				float32validator.AtLeast(3),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float32Value(1)),
					"Invalid Attribute Value",
					"Attribute test[Value(1.000000)] value must be at least 3.000000, got: 1.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float32Value(2)),
					"Invalid Attribute Value",
					"Attribute test[Value(2.000000)] value must be at least 3.000000, got: 2.000000",
				),
			},
		},
		"Set elements invalid for multiple validator": {
			val: types.SetValueMust(
				types.Float32Type,
				[]attr.Value{
					types.Float32Value(1),
					types.Float32Value(2),
				},
			),
			elementValidators: []validator.Float32{
				float32validator.AtLeast(3),
				float32validator.AtLeast(4),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float32Value(1)),
					"Invalid Attribute Value",
					"Attribute test[Value(1.000000)] value must be at least 3.000000, got: 1.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float32Value(1)),
					"Invalid Attribute Value",
					"Attribute test[Value(1.000000)] value must be at least 4.000000, got: 1.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float32Value(2)),
					"Invalid Attribute Value",
					"Attribute test[Value(2.000000)] value must be at least 3.000000, got: 2.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float32Value(2)),
					"Invalid Attribute Value",
					"Attribute test[Value(2.000000)] value must be at least 4.000000, got: 2.000000",
				),
			},
		},
		"Set elements valid": {
			val: types.SetValueMust(
				types.Float32Type,
				[]attr.Value{
					types.Float32Value(1),
					types.Float32Value(2),
				},
			),
			elementValidators: []validator.Float32{
				float32validator.AtLeast(1),
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
			setvalidator.ValueFloat32sAre(testCase.elementValidators...).ValidateSet(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
