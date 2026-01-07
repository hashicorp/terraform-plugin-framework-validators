// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package float32validator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float32validator"
)

func TestAtLeastValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Float32
		min         float32
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Float32": {
			val: types.Float32Unknown(),
			min: 0.90,
		},
		"null Float32": {
			val: types.Float32Null(),
			min: 0.90,
		},
		"valid integer as Float32": {
			val: types.Float32Value(2),
			min: 0.90,
		},
		"valid float as Float32": {
			val: types.Float32Value(2.2),
			min: 0.90,
		},
		"valid float as Float32 min": {
			val: types.Float32Value(0.9),
			min: 0.90,
		},
		"too small float as Float32": {
			val:         types.Float32Value(-1.1111),
			min:         0.90,
			expectError: true,
		},
	}

	for name, test := range tests {

		t.Run(fmt.Sprintf("ValidateFloat32 - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.Float32Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Float32Response{}
			float32validator.AtLeast(test.min).ValidateFloat32(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterFloat32 - %s", name), func(t *testing.T) {
			t.Parallel()
			request := function.Float32ParameterValidatorRequest{
				Value: test.val,
			}
			response := function.Float32ParameterValidatorResponse{}
			float32validator.AtLeast(test.min).ValidateParameterFloat32(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
