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

func TestAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Int64
		max         int64
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Int64": {
			val: types.Int64Unknown(),
			max: 2,
		},
		"null Int64": {
			val: types.Int64Null(),
			max: 2,
		},
		"valid integer as Int64": {
			val: types.Int64Value(1),
			max: 2,
		},
		"valid integer as Int64 min": {
			val: types.Int64Value(2),
			max: 2,
		},
		"too large integer as Int64": {
			val:         types.Int64Value(4),
			max:         2,
			expectError: true,
		},
		// Unknown value will be > 2
		"unknown lower bound exclusive - invalid less than bound": {
			val:         types.Int64Unknown().RefineWithLowerBound(2, false),
			max:         1,
			expectError: true,
		},
		"unknown lower bound exclusive - invalid matches bound": {
			val:         types.Int64Unknown().RefineWithLowerBound(2, false),
			max:         2,
			expectError: true,
		},
		"unknown lower bound exclusive - valid greater than bound": {
			val: types.Int64Unknown().RefineWithLowerBound(2, false),
			max: 3,
		},
		// Unknown value will be >= 2
		"unknown lower bound inclusive - invalid less than bound": {
			val:         types.Int64Unknown().RefineWithLowerBound(2, true),
			max:         1,
			expectError: true,
		},
		"unknown lower bound inclusive - valid matches bound": {
			val: types.Int64Unknown().RefineWithLowerBound(2, true),
			max: 2,
		},
		"unknown lower bound inclusive - valid greater than bound": {
			val: types.Int64Unknown().RefineWithLowerBound(2, true),
			max: 3,
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
			int64validator.AtMost(test.max).ValidateInt64(context.TODO(), request, &response)

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
			int64validator.AtMost(test.max).ValidateParameterInt64(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
