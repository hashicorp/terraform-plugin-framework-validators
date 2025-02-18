// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float64validator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
)

func TestAllValidatorValidateFloat64(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val        types.Float64
		validators []validator.Float64
		expected   diag.Diagnostics
	}
	tests := map[string]testCase{
		"invalid": {
			val: types.Float64Value(1.2),
			validators: []validator.Float64{
				float64validator.AtLeast(3),
				float64validator.AtLeast(5),
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value",
					"Attribute test value must be at least 3.000000, got: 1.200000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value",
					"Attribute test value must be at least 5.000000, got: 1.200000",
				),
			},
		},
		"valid": {
			val: types.Float64Value(1.2),
			validators: []validator.Float64{
				float64validator.AtLeast(0),
				float64validator.AtLeast(1),
			},
			expected: nil,
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.Float64Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Float64Response{}
			float64validator.All(test.validators...).ValidateFloat64(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
