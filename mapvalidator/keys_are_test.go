package mapvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestKeysAreValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val               attr.Value
		keysAreValidators []tfsdk.AttributeValidator
		expectError       bool
	}
	tests := map[string]testCase{
		"not Map": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{},
			),
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
		"Map key invalid": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one": types.StringValue("first"),
					"two": types.StringValue("second"),
				},
			),
			keysAreValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(4),
			},
			expectError: true,
		},
		"Map key invalid for second validator": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one": types.StringValue("first"),
					"two": types.StringValue("second"),
				},
			),
			keysAreValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(2),
				stringvalidator.LengthAtLeast(6),
			},
			expectError: true,
		},
		"Map keys wrong type for validator": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one": types.StringValue("first"),
					"two": types.StringValue("second"),
				},
			),
			keysAreValidators: []tfsdk.AttributeValidator{
				int64validator.AtLeast(6),
			},
			expectError: true,
		},
		"Map keys valid": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one": types.StringValue("first"),
					"two": types.StringValue("second"),
				},
			),
			keysAreValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(3),
			},
			expectError: false,
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
			KeysAre(test.keysAreValidators...).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
