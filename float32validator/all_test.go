// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float32validator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float32validator"
)

func TestAllValidatorValidateFloat32(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val        types.Float32
		validators []validator.Float32
		expected   diag.Diagnostics
	}
	tests := map[string]testCase{
		"invalid": {
			val: types.Float32Value(1.2),
			validators: []validator.Float32{
				float32validator.AtLeast(3),
				float32validator.AtLeast(5),
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
			val: types.Float32Value(1.2),
			validators: []validator.Float32{
				float32validator.AtLeast(0),
				float32validator.AtLeast(1),
			},
			expected: nil,
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.Float32Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Float32Response{}
			float32validator.All(test.validators...).ValidateFloat32(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
