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

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/testvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
)

func TestAnyWithAllWarningsValidatorValidateList(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val        types.List
		validators []validator.List
		expected   diag.Diagnostics
	}
	tests := map[string]testCase{
		"invalid": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			validators: []validator.List{
				listvalidator.SizeAtLeast(3),
				listvalidator.SizeAtLeast(5),
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value",
					"Attribute test list must contain at least 3 elements, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value",
					"Attribute test list must contain at least 5 elements, got: 2",
				),
			},
		},
		"valid": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			validators: []validator.List{
				listvalidator.SizeAtLeast(5),
				listvalidator.SizeAtLeast(2),
			},
			expected: diag.Diagnostics{},
		},
		"valid with warning": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			validators: []validator.List{
				listvalidator.All(listvalidator.SizeAtLeast(5), testvalidator.WarningList("failing warning summary", "failing warning details")),
				listvalidator.All(listvalidator.SizeAtLeast(2), testvalidator.WarningList("passing warning summary", "passing warning details")),
			},
			expected: diag.Diagnostics{
				diag.NewWarningDiagnostic("failing warning summary", "failing warning details"),
				diag.NewWarningDiagnostic("passing warning summary", "passing warning details"),
			},
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.ListRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.ListResponse{}
			listvalidator.AnyWithAllWarnings(test.validators...).ValidateList(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
