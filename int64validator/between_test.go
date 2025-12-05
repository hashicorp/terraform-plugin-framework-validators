// Copyright IBM Corp. 2022, 2025
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

func TestBetweenValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Int64
		min         int64
		max         int64
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Int64": {
			val: types.Int64Unknown(),
			min: 1,
			max: 3,
		},
		"null Int64": {
			val: types.Int64Null(),
			min: 1,
			max: 3,
		},
		"valid integer as Int64": {
			val: types.Int64Value(2),
			min: 1,
			max: 3,
		},
		"valid integer as Int64 min": {
			val: types.Int64Value(1),
			min: 1,
			max: 3,
		},
		"valid integer as Int64 max": {
			val: types.Int64Value(3),
			min: 1,
			max: 3,
		},
		"too small integer as Int64": {
			val:         types.Int64Value(-1),
			min:         1,
			max:         3,
			expectError: true,
		},
		"too large integer as Int64": {
			val:         types.Int64Value(42),
			min:         1,
			max:         3,
			expectError: true,
		},
		"invalid validator usage - minVal > maxVal": {
			val:         types.Int64Value(2),
			min:         3,
			max:         1,
			expectError: true,
		},
	}

	for name, test := range tests {

		t.Run(fmt.Sprintf("ValidateInt64 - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.Int64Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Int64Response{}
			int64validator.Between(test.min, test.max).ValidateInt64(context.TODO(), request, &response)

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
				ArgumentPosition: 0,
				Value:            test.val,
			}
			response := function.Int64ParameterValidatorResponse{}
			int64validator.Between(test.min, test.max).ValidateParameterInt64(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
