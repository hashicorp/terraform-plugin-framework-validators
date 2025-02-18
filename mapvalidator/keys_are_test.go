// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestKeysAreValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val               types.Map
		keysAreValidators []validator.String
		expectErrorsCount int
	}
	tests := map[string]testCase{
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
			keysAreValidators: []validator.String{
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
			keysAreValidators: []validator.String{
				stringvalidator.LengthAtLeast(2),
				stringvalidator.LengthAtLeast(6),
			},
			expectErrorsCount: 2,
		},
		"Map keys for invalid multiple validators": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one": types.StringValue("first"),
				},
			),
			keysAreValidators: []validator.String{
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
			keysAreValidators: []validator.String{
				stringvalidator.LengthAtLeast(3),
			},
			expectErrorsCount: 0,
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.MapRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.MapResponse{}
			KeysAre(test.keysAreValidators...).ValidateMap(context.TODO(), request, &response)

			if response.Diagnostics.ErrorsCount() != test.expectErrorsCount {
				t.Fatalf("expected %d errors, but got %d: %s", test.expectErrorsCount, response.Diagnostics.ErrorsCount(), response.Diagnostics)
			}
		})
	}
}
