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

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/testvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
)

func TestAnyValidatorValidateSet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val        types.Set
		validators []validator.Set
		expected   diag.Diagnostics
	}
	tests := map[string]testCase{
		"invalid": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			validators: []validator.Set{
				setvalidator.SizeAtLeast(3),
				setvalidator.SizeAtLeast(5),
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value",
					"Attribute test set must contain at least 3 elements, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value",
					"Attribute test set must contain at least 5 elements, got: 2",
				),
			},
		},
		"valid": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			validators: []validator.Set{
				setvalidator.SizeAtLeast(4),
				setvalidator.SizeAtLeast(2),
			},
			expected: diag.Diagnostics{},
		},
		"valid with warning": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			validators: []validator.Set{
				setvalidator.All(setvalidator.SizeAtLeast(5), testvalidator.WarningSet("failing warning summary", "failing warning details")),
				setvalidator.All(setvalidator.SizeAtLeast(2), testvalidator.WarningSet("passing warning summary", "passing warning details")),
			},
			expected: diag.Diagnostics{
				diag.NewWarningDiagnostic("passing warning summary", "passing warning details"),
			},
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.SetRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.SetResponse{}
			setvalidator.Any(test.validators...).ValidateSet(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
