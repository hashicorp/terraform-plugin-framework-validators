package mapvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSizeBetweenValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		min         int
		max         int
		expectError bool
	}
	tests := map[string]testCase{
		"not a Map": {
			val:         types.BoolValue(true),
			expectError: true,
		},
		"Map unknown": {
			val: types.MapUnknown(
				types.StringType,
			),
			expectError: false,
		},
		"Map null": {
			val: types.MapNull(
				types.StringType,
			),
			expectError: false,
		},
		"Map size greater than min": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one": types.StringValue("first"),
					"two": types.StringValue("second"),
				},
			),
			min:         1,
			max:         3,
			expectError: false,
		},
		"Map size equal to min": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one": types.StringValue("first"),
				},
			),
			min:         1,
			max:         3,
			expectError: false,
		},
		"Map size less than max": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one": types.StringValue("first"),
					"two": types.StringValue("second"),
				},
			),
			min:         1,
			max:         3,
			expectError: false,
		},
		"Map size equal to max": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one":   types.StringValue("first"),
					"two":   types.StringValue("second"),
					"three": types.StringValue("third"),
				},
			),
			min:         1,
			max:         3,
			expectError: false,
		},
		"Map size less than min": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{},
			),
			min:         1,
			max:         3,
			expectError: true,
		},
		"Map size greater than max": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one":   types.StringValue("first"),
					"two":   types.StringValue("second"),
					"three": types.StringValue("third"),
					"four":  types.StringValue("fourth"),
				},
			),
			min:         1,
			max:         3,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := tfsdk.ValidateAttributeRequest{
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
				AttributeConfig:         test.val,
			}
			response := tfsdk.ValidateAttributeResponse{}
			SizeBetween(test.min, test.max).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
