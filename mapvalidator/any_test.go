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

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/testvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
)

func TestAnyValidatorValidateMap(t *testing.T) {
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
				mapvalidator.SizeAtLeast(4),
				mapvalidator.SizeAtLeast(2),
			},
			expected: diag.Diagnostics{},
		},
		"valid with warning": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("first"),
					"key2": types.StringValue("second"),
				},
			),
			validators: []validator.Map{
				mapvalidator.All(mapvalidator.SizeAtLeast(5), testvalidator.WarningMap("failing warning summary", "failing warning details")),
				mapvalidator.All(mapvalidator.SizeAtLeast(2), testvalidator.WarningMap("passing warning summary", "passing warning details")),
			},
			expected: diag.Diagnostics{
				diag.NewWarningDiagnostic("passing warning summary", "passing warning details"),
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.MapRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.MapResponse{}
			mapvalidator.Any(test.validators...).ValidateMap(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
