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
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/testvalidator"
)

func TestAnyWithAllWarningsValidatorValidateFloat64(t *testing.T) {
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
			val: types.Float64Value(4),
			validators: []validator.Float64{
				float64validator.AtLeast(5),
				float64validator.AtLeast(3),
			},
			expected: diag.Diagnostics{},
		},
		"valid with warning": {
			val: types.Float64Value(4),
			validators: []validator.Float64{
				float64validator.All(float64validator.AtLeast(5), testvalidator.WarningFloat64("failing warning summary", "failing warning details")),
				float64validator.All(float64validator.AtLeast(2), testvalidator.WarningFloat64("passing warning summary", "passing warning details")),
			},
			expected: diag.Diagnostics{
				diag.NewWarningDiagnostic("failing warning summary", "failing warning details"),
				diag.NewWarningDiagnostic("passing warning summary", "passing warning details"),
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := validator.Float64Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Float64Response{}
			float64validator.AnyWithAllWarnings(test.validators...).ValidateFloat64(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
