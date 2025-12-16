// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package mapvalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
)

func TestAllValidatorValidateMap(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val        types.Map
		validators []validator.Map
		expected   diag.Diagnostics
	}
	tests := map[string]testCase{
		"invalid": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("first"),
					"key2": types.StringValue("second"),
				},
			),
			validators: []validator.Map{
				mapvalidator.SizeAtLeast(3),
				mapvalidator.SizeAtLeast(5),
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value",
					"Attribute test map must contain at least 3 elements, got: 2",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value",
					"Attribute test map must contain at least 5 elements, got: 2",
				),
			},
		},
		"valid": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("first"),
					"key2": types.StringValue("second"),
				},
			),
			validators: []validator.Map{
				mapvalidator.SizeAtLeast(0),
				mapvalidator.SizeAtLeast(1),
			},
			expected: nil,
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.MapRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.MapResponse{}
			mapvalidator.All(test.validators...).ValidateMap(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
