// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/testvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestAnyValidatorValidateString(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val        types.String
		validators []validator.String
		expected   diag.Diagnostics
	}
	tests := map[string]testCase{
		"invalid": {
			val: types.StringValue("one"),
			validators: []validator.String{
				stringvalidator.LengthAtLeast(4),
				stringvalidator.LengthAtLeast(5),
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value Length",
					"Attribute test string length must be at least 4, got: 3",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value Length",
					"Attribute test string length must be at least 5, got: 3",
				),
			},
		},
		"valid": {
			val: types.StringValue("test"),
			validators: []validator.String{
				stringvalidator.LengthAtLeast(5),
				stringvalidator.LengthAtLeast(3),
			},
			expected: diag.Diagnostics{},
		},
		"valid with warning": {
			val: types.StringValue("test"),
			validators: []validator.String{
				stringvalidator.All(stringvalidator.LengthAtLeast(5), testvalidator.WarningString("failing warning summary", "failing warning details")),
				stringvalidator.All(stringvalidator.LengthAtLeast(2), testvalidator.WarningString("passing warning summary", "passing warning details")),
			},
			expected: diag.Diagnostics{
				diag.NewWarningDiagnostic("passing warning summary", "passing warning details"),
			},
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.StringResponse{}
			stringvalidator.Any(test.validators...).ValidateString(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
