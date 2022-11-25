package setvalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
)

func TestValueFloat64sAreValidatorValidateSet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Set
		elementValidators   []validator.Float64
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.SetValueMust(
				types.Float64Type,
				[]attr.Value{
					types.Float64Value(1),
					types.Float64Value(2),
				},
			),
		},
		"Set unknown": {
			val: types.SetUnknown(
				types.Float64Type,
			),
			elementValidators: []validator.Float64{
				float64validator.AtLeast(1),
			},
		},
		"Set null": {
			val: types.SetNull(
				types.Float64Type,
			),
			elementValidators: []validator.Float64{
				float64validator.AtLeast(1),
			},
		},
		"Set elements invalid": {
			val: types.SetValueMust(
				types.Float64Type,
				[]attr.Value{
					types.Float64Value(1),
					types.Float64Value(2),
				},
			),
			elementValidators: []validator.Float64{
				float64validator.AtLeast(3),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float64Value(1)),
					"Invalid Attribute Value",
					"Attribute test[Value(1.000000)] value must be at least 3.000000, got: 1.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float64Value(2)),
					"Invalid Attribute Value",
					"Attribute test[Value(2.000000)] value must be at least 3.000000, got: 2.000000",
				),
			},
		},
		"Set elements invalid for multiple validator": {
			val: types.SetValueMust(
				types.Float64Type,
				[]attr.Value{
					types.Float64Value(1),
					types.Float64Value(2),
				},
			),
			elementValidators: []validator.Float64{
				float64validator.AtLeast(3),
				float64validator.AtLeast(4),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float64Value(1)),
					"Invalid Attribute Value",
					"Attribute test[Value(1.000000)] value must be at least 3.000000, got: 1.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float64Value(1)),
					"Invalid Attribute Value",
					"Attribute test[Value(1.000000)] value must be at least 4.000000, got: 1.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float64Value(2)),
					"Invalid Attribute Value",
					"Attribute test[Value(2.000000)] value must be at least 3.000000, got: 2.000000",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtSetValue(types.Float64Value(2)),
					"Invalid Attribute Value",
					"Attribute test[Value(2.000000)] value must be at least 4.000000, got: 2.000000",
				),
			},
		},
		"Set elements valid": {
			val: types.SetValueMust(
				types.Float64Type,
				[]attr.Value{
					types.Float64Value(1),
					types.Float64Value(2),
				},
			),
			elementValidators: []validator.Float64{
				float64validator.AtLeast(1),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.SetRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.val,
			}
			response := validator.SetResponse{}
			setvalidator.ValueFloat64sAre(testCase.elementValidators...).ValidateSet(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
