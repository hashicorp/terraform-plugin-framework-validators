// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSizeAtLeastValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Map
		min         int
		expectError bool
	}
	tests := map[string]testCase{
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
			expectError: false,
		},
		"Map size less than min": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{},
			),
			min:         1,
			expectError: true,
		},
		// Unknown value will have <= 2 elements
		"unknown length upper bound - valid less than bound": {
			val: types.MapUnknown(types.StringType).RefineWithLengthUpperBound(2),
			min: 1,
		},
		"unknown length upper bound - valid matches bound": {
			val: types.MapUnknown(types.StringType).RefineWithLengthUpperBound(2),
			min: 2,
		},
		"unknown length upper bound - invalid greater than bound": {
			val:         types.MapUnknown(types.StringType).RefineWithLengthUpperBound(2),
			min:         3,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(fmt.Sprintf("ValidateMap - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.MapRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.MapResponse{}
			SizeAtLeast(test.min).ValidateMap(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterMap - %s", name), func(t *testing.T) {
			t.Parallel()
			request := function.MapParameterValidatorRequest{
				ArgumentPosition: 0,
				Value:            test.val,
			}
			response := function.MapParameterValidatorResponse{}
			SizeAtLeast(test.min).ValidateParameterMap(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
