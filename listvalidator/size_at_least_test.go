// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator

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
		val         types.List
		min         int
		expectError bool
	}
	tests := map[string]testCase{
		"List unknown": {
			val: types.ListUnknown(
				types.StringType,
			),
			expectError: false,
		},
		"List null": {
			val: types.ListNull(
				types.StringType,
			),
			expectError: false,
		},
		"List size greater than min": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			min:         1,
			expectError: false,
		},
		"List size equal to min": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
				},
			),
			min:         1,
			expectError: false,
		},
		"List size less than min": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{},
			),
			min:         1,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(fmt.Sprintf("ValidateList - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.ListRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.ListResponse{}
			SizeAtLeast(test.min).ValidateList(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterList - %s", name), func(t *testing.T) {
			t.Parallel()
			request := function.ListParameterValidatorRequest{
				ArgumentPosition: 0,
				Value:            test.val,
			}
			response := function.ListParameterValidatorResponse{}
			SizeAtLeast(test.min).ValidateParameterList(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
