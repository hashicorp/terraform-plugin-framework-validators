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

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
)

func TestValueSetsAreValidatorValidateSet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Set
		elementValidators   []validator.Set
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.SetValueMust(
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
		"Set unknown": {
			val: types.SetUnknown(
				types.StringType,
			),
			elementValidators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
		"Set null": {
			val: types.SetNull(
				types.StringType,
			),
			elementValidators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
		"Set elements invalid": {
			val: types.SetValueMust(
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
					path.Root("test").AtSetValue(types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("first"),
							types.StringValue("second"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value([\"first\",\"second\"])] set must contain at least 3 elements, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("third"),
							types.StringValue("fourth"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value([\"third\",\"fourth\"])] set must contain at least 3 elements, got: 2",
				),
			},
		},
		"Set elements invalid for multiple validator": {
			val: types.SetValueMust(
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
					path.Root("test").AtSetValue(types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("first"),
							types.StringValue("second"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value([\"first\",\"second\"])] set must contain at least 3 elements, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("first"),
							types.StringValue("second"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value([\"first\",\"second\"])] set must contain at least 4 elements, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("third"),
							types.StringValue("fourth"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value([\"third\",\"fourth\"])] set must contain at least 3 elements, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("third"),
							types.StringValue("fourth"),
						},
					)),
					"Invalid Attribute Value",
					"Attribute test[Value([\"third\",\"fourth\"])] set must contain at least 4 elements, got: 2",
				),
			},
		},
		"Set elements valid": {
			val: types.SetValueMust(
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
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.SetRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.val,
			}
			response := validator.SetResponse{}
			setvalidator.ValueSetsAre(testCase.elementValidators...).ValidateSet(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
