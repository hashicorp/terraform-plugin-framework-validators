// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package setvalidator

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
		val         types.Set
		min         int
		expectError bool
	}
	tests := map[string]testCase{
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
			expectError: false,
		},
		"Set size less than min": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{},
			),
			min:         1,
			expectError: true,
		},
		// Unknown value will have <= 2 elements
		"unknown length upper bound - valid less than bound": {
			val: types.SetUnknown(types.StringType).RefineWithLengthUpperBound(2),
			min: 1,
		},
		"unknown length upper bound - valid matches bound": {
			val: types.SetUnknown(types.StringType).RefineWithLengthUpperBound(2),
			min: 2,
		},
		"unknown length upper bound - invalid greater than bound": {
			val:         types.SetUnknown(types.StringType).RefineWithLengthUpperBound(2),
			min:         3,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(fmt.Sprintf("ValidateSet - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.SetRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.SetResponse{}
			SizeAtLeast(test.min).ValidateSet(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterSet - %s", name), func(t *testing.T) {
			t.Parallel()
			request := function.SetParameterValidatorRequest{
				ArgumentPosition: 0,
				Value:            test.val,
			}
			response := function.SetParameterValidatorResponse{}
			SizeAtLeast(test.min).ValidateParameterSet(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
