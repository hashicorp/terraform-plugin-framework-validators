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

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
)

func TestValueInt32sAreValidatorValidateList(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.List
		elementValidators   []validator.Int32
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.ListValueMust(
				types.Int32Type,
				[]attr.Value{
					types.Int32Value(1),
					types.Int32Value(2),
				},
			),
		},
		"List unknown": {
			val: types.ListUnknown(
				types.Int32Type,
			),
			elementValidators: []validator.Int32{
				int32validator.AtLeast(1),
			},
		},
		"List null": {
			val: types.ListNull(
				types.Int32Type,
			),
			elementValidators: []validator.Int32{
				int32validator.AtLeast(1),
			},
		},
		"List elements invalid": {
			val: types.ListValueMust(
				types.Int32Type,
				[]attr.Value{
					types.Int32Value(1),
					types.Int32Value(2),
				},
			),
			elementValidators: []validator.Int32{
				int32validator.AtLeast(2),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value",
					"Attribute test[0] value must be at least 2, got: 1",
				),
			},
		},
		"List elements invalid for multiple validator": {
			val: types.ListValueMust(
				types.Int32Type,
				[]attr.Value{
					types.Int32Value(1),
					types.Int32Value(2),
				},
			),
			elementValidators: []validator.Int32{
				int32validator.AtLeast(3),
				int32validator.AtLeast(4),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value",
					"Attribute test[0] value must be at least 3, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value",
					"Attribute test[0] value must be at least 4, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(1),
					"Invalid Attribute Value",
					"Attribute test[1] value must be at least 3, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(1),
					"Invalid Attribute Value",
					"Attribute test[1] value must be at least 4, got: 2",
				),
			},
		},
		"List elements valid": {
			val: types.ListValueMust(
				types.Int32Type,
				[]attr.Value{
					types.Int32Value(1),
					types.Int32Value(2),
				},
			),
			elementValidators: []validator.Int32{
				int32validator.AtLeast(1),
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
			listvalidator.ValueInt32sAre(testCase.elementValidators...).ValidateList(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
