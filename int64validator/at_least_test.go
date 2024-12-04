// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int64validator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
)

func TestAtLeastValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Int64
		min         int64
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Int64": {
			val: types.Int64Unknown(),
			min: 1,
		},
		"null Int64": {
			val: types.Int64Null(),
			min: 1,
		},
		"valid integer as Int64": {
			val: types.Int64Value(2),
			min: 1,
		},
		"valid integer as Int64 min": {
			val: types.Int64Value(1),
			min: 1,
		},
		"too small integer as Int64": {
			val:         types.Int64Value(-1),
			min:         1,
			expectError: true,
		},
		// Unknown value will be < 2
		"unknown upper bound exclusive - valid less than bound": {
			val: types.Int64Unknown().RefineWithUpperBound(2, false),
			min: 1,
		},
		"unknown upper bound exclusive - invalid matches bound": {
			val:         types.Int64Unknown().RefineWithUpperBound(2, false),
			min:         2,
			expectError: true,
		},
		"unknown upper bound exclusive - invalid greater than bound": {
			val:         types.Int64Unknown().RefineWithUpperBound(2, false),
			min:         3,
			expectError: true,
		},
		// Unknown value will be <= 2
		"unknown upper bound inclusive - valid less than bound": {
			val: types.Int64Unknown().RefineWithUpperBound(2, true),
			min: 1,
		},
		"unknown upper bound inclusive - valid matches bound": {
			val: types.Int64Unknown().RefineWithUpperBound(2, true),
			min: 2,
		},
		"unknown upper bound inclusive - invalid greater than bound": {
			val:         types.Int64Unknown().RefineWithUpperBound(2, true),
			min:         3,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(fmt.Sprintf("ValidateInt64 - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.Int64Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Int64Response{}
			int64validator.AtLeast(test.min).ValidateInt64(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterInt64 - %s", name), func(t *testing.T) {
			t.Parallel()
			request := function.Int64ParameterValidatorRequest{
				Value: test.val,
			}
			response := function.Int64ParameterValidatorResponse{}
			int64validator.AtLeast(test.min).ValidateParameterInt64(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
