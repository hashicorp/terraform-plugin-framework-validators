// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int32validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
)

func TestBetweenValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Int32
		min         int32
		max         int32
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Int32": {
			val: types.Int32Unknown(),
			min: 1,
			max: 3,
		},
		"null Int32": {
			val: types.Int32Null(),
			min: 1,
			max: 3,
		},
		"valid integer as Int32": {
			val: types.Int32Value(2),
			min: 1,
			max: 3,
		},
		"valid integer as Int32 min": {
			val: types.Int32Value(1),
			min: 1,
			max: 3,
		},
		"valid integer as Int32 max": {
			val: types.Int32Value(3),
			min: 1,
			max: 3,
		},
		"too small integer as Int32": {
			val:         types.Int32Value(-1),
			min:         1,
			max:         3,
			expectError: true,
		},
		"too large integer as Int32": {
			val:         types.Int32Value(42),
			min:         1,
			max:         3,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.Int32Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Int32Response{}
			int32validator.Between(test.min, test.max).ValidateInt32(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
