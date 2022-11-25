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

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
)

func TestValueInt64sAreValidatorValidateMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Map
		elementValidators   []validator.Int64
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.MapValueMust(
				types.Int64Type,
				map[string]attr.Value{
					"key1": types.Int64Value(1),
					"key2": types.Int64Value(2),
				},
			),
		},
		"Map unknown": {
			val: types.MapUnknown(
				types.Int64Type,
			),
			elementValidators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"Map null": {
			val: types.MapNull(
				types.Int64Type,
			),
			elementValidators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"Map elements invalid": {
			val: types.MapValueMust(
				types.Int64Type,
				map[string]attr.Value{
					"key1": types.Int64Value(1),
					// Map ordering is random in Go, avoid multiple keys
				},
			),
			elementValidators: []validator.Int64{
				int64validator.AtLeast(3),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value",
					"Attribute test[\"key1\"] value must be at least 3, got: 1",
				),
			},
		},
		"Map elements invalid for multiple validator": {
			val: types.MapValueMust(
				types.Int64Type,
				map[string]attr.Value{
					"key1": types.Int64Value(1),
					// Map ordering is random in Go, avoid multiple keys
				},
			),
			elementValidators: []validator.Int64{
				int64validator.AtLeast(3),
				int64validator.AtLeast(4),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value",
					"Attribute test[\"key1\"] value must be at least 3, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value",
					"Attribute test[\"key1\"] value must be at least 4, got: 1",
				),
			},
		},
		"Map elements valid": {
			val: types.MapValueMust(
				types.Int64Type,
				map[string]attr.Value{
					"key1": types.Int64Value(1),
					"key2": types.Int64Value(2),
				},
			),
			elementValidators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.MapRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.val,
			}
			response := validator.MapResponse{}
			mapvalidator.ValueInt64sAre(testCase.elementValidators...).ValidateMap(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
