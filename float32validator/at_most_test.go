// Copyright (c) HashiCorp, Inc.
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

func TestAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Float32
		max         float32
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Float32": {
			val: types.Float32Unknown(),
			max: 2.00,
		},
		"null Float32": {
			val: types.Float32Null(),
			max: 2.00,
		},
		"valid integer as Float32": {
			val: types.Float32Value(1),
			max: 2.00,
		},
		"valid float as Float32": {
			val: types.Float32Value(1.1),
			max: 2.00,
		},
		"valid float as Float32 max": {
			val: types.Float32Value(2.0),
			max: 2.00,
		},
		"too large float as Float32": {
			val:         types.Float32Value(3.0),
			max:         2.00,
			expectError: true,
		},
		// Unknown value will be > 2.1
		"unknown lower bound exclusive - invalid less than bound": {
			val:         types.Float32Unknown().RefineWithLowerBound(2.1, false),
			max:         2,
			expectError: true,
		},
		"unknown lower bound exclusive - invalid matches bound": {
			val:         types.Float32Unknown().RefineWithLowerBound(2.1, false),
			max:         2.1,
			expectError: true,
		},
		"unknown lower bound exclusive - valid greater than bound": {
			val: types.Float32Unknown().RefineWithLowerBound(2.1, false),
			max: 3,
		},
		// Unknown value will be >= 2.1
		"unknown lower bound inclusive - invalid less than bound": {
			val:         types.Float32Unknown().RefineWithLowerBound(2.1, true),
			max:         2,
			expectError: true,
		},
		"unknown lower bound inclusive - valid matches bound": {
			val: types.Float32Unknown().RefineWithLowerBound(2.1, true),
			max: 2.1,
		},
		"unknown lower bound inclusive - valid greater than bound": {
			val: types.Float32Unknown().RefineWithLowerBound(2.1, true),
			max: 3,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(fmt.Sprintf("ValidateFloat32 - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.Float32Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Float32Response{}
			float32validator.AtMost(test.max).ValidateFloat32(context.TODO(), request, &response)

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
			float32validator.AtMost(test.max).ValidateParameterFloat32(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
