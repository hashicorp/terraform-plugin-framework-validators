package objectvalidator_test

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
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
)

func TestAnyWithAllWarningsValidatorValidateObject(t *testing.T) {
	t.Parallel()

	testValue := types.ObjectValueMust(
		map[string]attr.Type{
			"testattr": types.StringType,
		},
		map[string]attr.Value{
			"testattr": types.StringValue("test"),
		},
	)

	type testCase struct {
		val        types.Object
		validators []validator.Object
		expected   diag.Diagnostics
	}
	tests := map[string]testCase{
		"invalid": {
			val: testValue,
			validators: []validator.Object{
				testvalidator.ObjectValidator{
					Diagnostics: diag.Diagnostics{
						diag.NewAttributeErrorDiagnostic(
							path.Root("test"),
							"Error Summary 1",
							"Error Detail 1",
						),
					},
				},
				testvalidator.ObjectValidator{
					Diagnostics: diag.Diagnostics{
						diag.NewAttributeErrorDiagnostic(
							path.Root("test"),
							"Error Summary 2",
							"Error Detail 2",
						),
					},
				},
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Error Summary 1",
					"Error Detail 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Error Summary 2",
					"Error Detail 2",
				),
			},
		},
		"valid": {
			val: testValue,
			validators: []validator.Object{
				testvalidator.ObjectValidator{
					Diagnostics: diag.Diagnostics{
						diag.NewAttributeErrorDiagnostic(
							path.Root("test"),
							"Error Summary",
							"Error Detail",
						),
					},
				},
				testvalidator.ObjectValidator{},
			},
			expected: diag.Diagnostics{},
		},
		"valid with warning": {
			val: testValue,
			validators: []validator.Object{
				testvalidator.ObjectValidator{
					Diagnostics: diag.Diagnostics{
						diag.NewAttributeWarningDiagnostic(
							path.Root("test"),
							"Failing Warning Summary",
							"Failing Warning Detail",
						),
						diag.NewAttributeErrorDiagnostic(
							path.Root("test"),
							"Failing Error Summary",
							"Failing Error Detail",
						),
					},
				},
				testvalidator.ObjectValidator{
					Diagnostics: diag.Diagnostics{
						diag.NewAttributeWarningDiagnostic(
							path.Root("test"),
							"Passing Warning Summary",
							"Passing Warning Detail",
						),
					},
				},
			},
			expected: diag.Diagnostics{
				diag.NewAttributeWarningDiagnostic(
					path.Root("test"),
					"Failing Warning Summary",
					"Failing Warning Detail",
				),
				diag.NewAttributeWarningDiagnostic(
					path.Root("test"),
					"Passing Warning Summary",
					"Passing Warning Detail",
				),
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.ObjectRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.ObjectResponse{}
			objectvalidator.AnyWithAllWarnings(test.validators...).ValidateObject(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, test.expected); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
