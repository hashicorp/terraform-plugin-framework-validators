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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestValueStringsAreValidatorValidateSet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Set
		elementValidators   []validator.String
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
		},
		"Set unknown": {
			val: types.SetUnknown(
				types.StringType,
			),
			elementValidators: []validator.String{
				stringvalidator.LengthAtLeast(6),
			},
		},
		"Set null": {
			val: types.SetNull(
				types.StringType,
			),
			elementValidators: []validator.String{
				stringvalidator.LengthAtLeast(6),
			},
		},
		"Set elements invalid": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			elementValidators: []validator.String{
				stringvalidator.LengthAtLeast(7),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.StringValue("first")),
					"Invalid Attribute Value Length",
					"Attribute test[Value(\"first\")] string length must be at least 7, got: 5",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.StringValue("second")),
					"Invalid Attribute Value Length",
					"Attribute test[Value(\"second\")] string length must be at least 7, got: 6",
				),
			},
		},
		"Set elements invalid for multiple validator": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			elementValidators: []validator.String{
				stringvalidator.LengthAtLeast(7),
				stringvalidator.LengthAtLeast(8),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.StringValue("first")),
					"Invalid Attribute Value Length",
					"Attribute test[Value(\"first\")] string length must be at least 7, got: 5",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.StringValue("first")),
					"Invalid Attribute Value Length",
					"Attribute test[Value(\"first\")] string length must be at least 8, got: 5",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.StringValue("second")),
					"Invalid Attribute Value Length",
					"Attribute test[Value(\"second\")] string length must be at least 7, got: 6",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.StringValue("second")),
					"Invalid Attribute Value Length",
					"Attribute test[Value(\"second\")] string length must be at least 8, got: 6",
				),
			},
		},
		"Set elements valid": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			elementValidators: []validator.String{
				stringvalidator.LengthAtLeast(5),
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.SetRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.val,
			}
			response := validator.SetResponse{}
			setvalidator.ValueStringsAre(testCase.elementValidators...).ValidateSet(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
