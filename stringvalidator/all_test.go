package stringvalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestAllValidatorValidateString(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val        types.String
		validators []validator.String
		expected   diag.Diagnostics
	}
	tests := map[string]testCase{
		"invalid": {
			val: types.StringValue("test"),
			validators: []validator.String{
				stringvalidator.LengthAtLeast(5),
				stringvalidator.LengthAtLeast(6),
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value Length",
					"Attribute test string length must be at least 5, got: 4",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Value Length",
					"Attribute test string length must be at least 6, got: 4",
				),
			},
		},
		"valid": {
			val: types.StringValue("test"),
			validators: []validator.String{
				stringvalidator.LengthAtLeast(0),
				stringvalidator.LengthAtLeast(1),
			},
			expected: nil,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.StringResponse{}
			stringvalidator.All(test.validators...).ValidateString(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
