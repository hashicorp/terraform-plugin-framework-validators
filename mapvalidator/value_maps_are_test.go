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

func TestValueMapsAreValidatorValidateMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		val                 types.Map
		elementValidators   []validator.Map
		expectedDiagnostics diag.Diagnostics
	}{
		"no element validators": {
			val: types.MapValueMust(
				types.MapType{ElemType: types.StringType},
				map[string]attr.Value{
					"key1": types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("first"),
							"key2": types.StringValue("second"),
						},
					),
					"key2": types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("third"),
							"key2": types.StringValue("fourth"),
						},
					),
				},
			),
		},
		"Map unknown": {
			val: types.MapUnknown(
				types.StringType,
			),
			elementValidators: []validator.Map{
				mapvalidator.SizeAtLeast(1),
			},
		},
		"Map null": {
			val: types.MapNull(
				types.StringType,
			),
			elementValidators: []validator.Map{
				mapvalidator.SizeAtLeast(1),
			},
		},
		"Map elements invalid": {
			val: types.MapValueMust(
				types.MapType{ElemType: types.StringType},
				map[string]attr.Value{
					"key1": types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("first"),
							// Map ordering is random in Go, avoid multiple keys
						},
					),
					// Map ordering is random in Go, avoid multiple keys
				},
			),
			elementValidators: []validator.Map{
				mapvalidator.SizeAtLeast(3),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value",
					"Attribute test[\"key1\"] map must contain at least 3 elements, got: 1",
				),
			},
		},
		"Map elements invalid for multiple validator": {
			val: types.MapValueMust(
				types.MapType{ElemType: types.StringType},
				map[string]attr.Value{
					"key1": types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("first"),
							// Map ordering is random in Go, avoid multiple keys
						},
					),
					// Map ordering is random in Go, avoid multiple keys
				},
			),
			elementValidators: []validator.Map{
				mapvalidator.SizeAtLeast(3),
				mapvalidator.SizeAtLeast(4),
			},
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value",
					"Attribute test[\"key1\"] map must contain at least 3 elements, got: 1",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test").AtMapKey("key1"),
					"Invalid Attribute Value",
					"Attribute test[\"key1\"] map must contain at least 4 elements, got: 1",
				),
			},
		},
		"Map elements valid": {
			val: types.MapValueMust(
				types.MapType{ElemType: types.StringType},
				map[string]attr.Value{
					"key1": types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("first"),
							"key2": types.StringValue("second"),
						},
					),
					"key2": types.MapValueMust(
						types.StringType,
						map[string]attr.Value{
							"key1": types.StringValue("third"),
							"key2": types.StringValue("fourth"),
						},
					),
				},
			),
			elementValidators: []validator.Map{
				mapvalidator.SizeAtLeast(1),
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
			mapvalidator.ValueMapsAre(testCase.elementValidators...).ValidateMap(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
