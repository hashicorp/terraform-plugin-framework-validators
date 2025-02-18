// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int32validator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/testvalidator"
)

func TestAnyValidatorValidateInt32(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val        types.Int32
		validators []validator.Int32
		expected   diag.Diagnostics
	}
	tests := map[string]testCase{
		"invalid": {
			val: types.Int32Value(1),
			validators: []validator.Int32{
				int32validator.AtLeast(3),
				int32validator.AtLeast(5),
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
			val: types.Int32Value(4),
			validators: []validator.Int32{
				int32validator.AtLeast(5),
				int32validator.AtLeast(3),
			},
			expected: diag.Diagnostics{},
		},
		"valid with warning": {
			val: types.Int32Value(4),
			validators: []validator.Int32{
				int32validator.All(int32validator.AtLeast(5), testvalidator.WarningInt32("failing warning summary", "failing warning details")),
				int32validator.All(int32validator.AtLeast(2), testvalidator.WarningInt32("passing warning summary", "passing warning details")),
			},
			expected: diag.Diagnostics{
				diag.NewWarningDiagnostic("passing warning summary", "passing warning details"),
			},
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.Int32Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Int32Response{}
			int32validator.Any(test.validators...).ValidateInt32(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
