// Copyright IBM Corp. 2022, 2025
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

func TestAtLeastValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Int32
		min         int32
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Int32": {
			val: types.Int32Unknown(),
			min: 1,
		},
		"null Int32": {
			val: types.Int32Null(),
			min: 1,
		},
		"valid integer as Int32": {
			val: types.Int32Value(2),
			min: 1,
		},
		"valid integer as Int32 min": {
			val: types.Int32Value(1),
			min: 1,
		},
		"too small integer as Int32": {
			val:         types.Int32Value(-1),
			min:         1,
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
			int32validator.AtLeast(test.min).ValidateInt32(context.TODO(), request, &response)

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
			int32validator.AtLeast(test.min).ValidateParameterInt32(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
