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
		expectErrorsCount int
	}
	tests := map[string]testCase{
		"not Map": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{},
			),
			expectErrorsCount: 1,
		},
		"Map unknown": {
			val: types.MapUnknown(
				types.StringType,
			),
			expectErrorsCount: 0,
		},
		"Map null": {
			val: types.MapNull(
				types.StringType,
			),
			expectErrorsCount: 0,
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
			expectErrorsCount: 2,
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
			expectErrorsCount: 2,
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
			expectErrorsCount: 1,
		},
		"Map keys for invalid multiple validators": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one": types.StringValue("first"),
				},
			),
			keysAreValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(5),
				stringvalidator.LengthAtLeast(6),
			},
			expectErrorsCount: 2,
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
			expectErrorsCount: 0,
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

			if response.Diagnostics.ErrorsCount() != test.expectErrorsCount {
				t.Fatalf("expected %d errors, but got %d: %s", test.expectErrorsCount, response.Diagnostics.ErrorsCount(), response.Diagnostics)
			}
		})
	}
}
