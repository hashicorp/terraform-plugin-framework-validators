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

	"github.com/hashicorp/terraform-plugin-framework-validators/float32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
)

func TestValueFloat32sAreValidatorValidateList(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.List
		elementValidators   []validator.Float32
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.ListValueMust(
				types.Float32Type,
				[]attr.Value{
					types.Float32Value(1),
					types.Float32Value(2),
				},
			),
		},
		"List unknown": {
			val: types.ListUnknown(
				types.Float32Type,
			),
			elementValidators: []validator.Float32{
				float32validator.AtLeast(1),
			},
		},
		"List null": {
			val: types.ListNull(
				types.Float32Type,
			),
			elementValidators: []validator.Float32{
				float32validator.AtLeast(1),
			},
		},
		"List elements invalid": {
			val: types.ListValueMust(
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

			request := validator.ListRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.val,
			}
			response := validator.ListResponse{}
			listvalidator.ValueFloat32sAre(testCase.elementValidators...).ValidateList(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
