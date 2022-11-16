package setvalidator

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
		"not a Set": {
			val:         types.BoolValue(true),
			expectError: true,
		},
		"Set unknown": {
			val: types.SetUnknown(
				types.StringType,
			),
			expectError: false,
		},
		"Set null": {
			val: types.SetNull(
				types.StringType,
			),
			expectError: false,
		},
		"Set size greater than min": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			min:         1,
			max:         3,
			expectError: false,
		},
		"Set size equal to min": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
				},
			),
			min:         1,
			max:         3,
			expectError: false,
		},
		"Set size less than max": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			min:         1,
			max:         3,
			expectError: false,
		},
		"Set size equal to max": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
					types.StringValue("third"),
				},
			),
			min:         1,
			max:         3,
			expectError: false,
		},
		"Set size less than min": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{},
			),
			min:         1,
			max:         3,
			expectError: true,
		},
		"Set size greater than max": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
					types.StringValue("third"),
					types.StringValue("fourth"),
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
