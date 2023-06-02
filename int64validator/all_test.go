// Copyright (c) HashiCorp, Inc.
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
)

func TestAllValidatorValidateInt64(t *testing.T) {
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
			val: types.Int64Value(1),
			validators: []validator.Int64{
				int64validator.AtLeast(0),
				int64validator.AtLeast(1),
			},
			expected: nil,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.Int64Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Int64Response{}
			int64validator.All(test.validators...).ValidateInt64(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
