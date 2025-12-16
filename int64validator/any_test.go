// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package int64validator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/testvalidator"
)

func TestAnyValidatorValidateInt64(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val        types.Int64
		validators []validator.Int64
		expected   diag.Diagnostics
	}
	tests := map[string]testCase{
		"invalid": {
			val: types.Int64Value(1),
			validators: []validator.Int64{
				int64validator.AtLeast(3),
				int64validator.AtLeast(5),
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value",
					"Attribute test value must be at least 3, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value",
					"Attribute test value must be at least 5, got: 1",
				),
			},
		},
		"valid": {
			val: types.Int64Value(4),
			validators: []validator.Int64{
				int64validator.AtLeast(5),
				int64validator.AtLeast(3),
			},
			expected: diag.Diagnostics{},
		},
		"valid with warning": {
			val: types.Int64Value(4),
			validators: []validator.Int64{
				int64validator.All(int64validator.AtLeast(5), testvalidator.WarningInt64("failing warning summary", "failing warning details")),
				int64validator.All(int64validator.AtLeast(2), testvalidator.WarningInt64("passing warning summary", "passing warning details")),
			},
			expected: diag.Diagnostics{
				diag.NewWarningDiagnostic("passing warning summary", "passing warning details"),
			},
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.Int64Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Int64Response{}
			int64validator.Any(test.validators...).ValidateInt64(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
