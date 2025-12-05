// Copyright IBM Corp. 2022, 2025
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

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
)

func TestValueSetsAreValidatorValidateList(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.List
		elementValidators   []validator.Set
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.ListValueMust(
				types.SetType{ElemType: types.StringType},
				[]attr.Value{
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("first"),
							types.StringValue("second"),
						},
					),
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("third"),
							types.StringValue("fourth"),
						},
					),
				},
			),
		},
		"List unknown": {
			val: types.ListUnknown(
				types.StringType,
			),
			elementValidators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
		"List null": {
			val: types.ListNull(
				types.StringType,
			),
			elementValidators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
		"List elements invalid": {
			val: types.ListValueMust(
				types.SetType{ElemType: types.StringType},
				[]attr.Value{
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("first"),
							types.StringValue("second"),
						},
					),
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("third"),
							types.StringValue("fourth"),
						},
					),
				},
			),
			elementValidators: []validator.Set{
				setvalidator.SizeAtLeast(3),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value",
					"Attribute test[0] set must contain at least 3 elements, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(1),
					"Invalid Attribute Value",
					"Attribute test[1] set must contain at least 3 elements, got: 2",
				),
			},
		},
		"List elements invalid for multiple validator": {
			val: types.ListValueMust(
				types.SetType{ElemType: types.StringType},
				[]attr.Value{
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("first"),
							types.StringValue("second"),
						},
					),
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("third"),
							types.StringValue("fourth"),
						},
					),
				},
			),
			elementValidators: []validator.Set{
				setvalidator.SizeAtLeast(3),
				setvalidator.SizeAtLeast(4),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value",
					"Attribute test[0] set must contain at least 3 elements, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(0),
					"Invalid Attribute Value",
					"Attribute test[0] set must contain at least 4 elements, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(1),
					"Invalid Attribute Value",
					"Attribute test[1] set must contain at least 3 elements, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtListIndex(1),
					"Invalid Attribute Value",
					"Attribute test[1] set must contain at least 4 elements, got: 2",
				),
			},
		},
		"List elements valid": {
			val: types.ListValueMust(
				types.SetType{ElemType: types.StringType},
				[]attr.Value{
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("first"),
							types.StringValue("second"),
						},
					),
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("third"),
							types.StringValue("fourth"),
						},
					),
				},
			),
			elementValidators: []validator.Set{
				setvalidator.SizeAtLeast(1),
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
			listvalidator.ValueSetsAre(testCase.elementValidators...).ValidateList(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
