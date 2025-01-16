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

func TestSizeAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Set
		max         int
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
		"Set size less than max": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
				},
			),
			max:         2,
			expectError: false,
		},
		"Set size equal to max": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			max:         2,
			expectError: false,
		},
		"Set size greater than max": {
			val: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
					types.StringValue("third"),
				},
			),
			max:         2,
			expectError: true,
		},
		// Unknown value will have >= 2 elements
		"unknown length lower bound - invalid less than bound": {
			val:         types.SetUnknown(types.StringType).RefineWithLengthLowerBound(2),
			max:         1,
			expectError: true,
		},
		"unknown length lower bound - valid matches bound": {
			val: types.SetUnknown(types.StringType).RefineWithLengthLowerBound(2),
			max: 2,
		},
		"unknown length lower bound - valid greater than bound": {
			val: types.SetUnknown(types.StringType).RefineWithLengthLowerBound(2),
			max: 3,
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
			SizeAtMost(test.max).ValidateSet(context.TODO(), request, &response)

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
			SizeAtMost(test.max).ValidateParameterSet(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
