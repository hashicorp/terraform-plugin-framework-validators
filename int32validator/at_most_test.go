// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int32validator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
)

func TestAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Int32
		max         int32
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Int32": {
			val: types.Int32Unknown(),
			max: 2,
		},
		"null Int32": {
			val: types.Int32Null(),
			max: 2,
		},
		"valid integer as Int32": {
			val: types.Int32Value(1),
			max: 2,
		},
		"valid integer as Int32 min": {
			val: types.Int32Value(2),
			max: 2,
		},
		"too large integer as Int32": {
			val:         types.Int32Value(4),
			max:         2,
			expectError: true,
		},
	}

	for name, test := range tests {

		t.Run(fmt.Sprintf("ValidateInt32 - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.Int32Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Int32Response{}
			int32validator.AtMost(test.max).ValidateInt32(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterInt32 - %s", name), func(t *testing.T) {
			t.Parallel()
			request := function.Int32ParameterValidatorRequest{
				Value: test.val,
			}
			response := function.Int32ParameterValidatorResponse{}
			int32validator.AtMost(test.max).ValidateParameterInt32(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
