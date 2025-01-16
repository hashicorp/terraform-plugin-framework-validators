// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float64validator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
)

func TestAtLeastValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Float64
		min         float64
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Float64": {
			val: types.Float64Unknown(),
			min: 0.90,
		},
		"null Float64": {
			val: types.Float64Null(),
			min: 0.90,
		},
		"valid integer as Float64": {
			val: types.Float64Value(2),
			min: 0.90,
		},
		"valid float as Float64": {
			val: types.Float64Value(2.2),
			min: 0.90,
		},
		"valid float as Float64 min": {
			val: types.Float64Value(0.9),
			min: 0.90,
		},
		"too small float as Float64": {
			val:         types.Float64Value(-1.1111),
			min:         0.90,
			expectError: true,
		},
		// Unknown value will be < 2.1
		"unknown upper bound exclusive - valid less than bound": {
			val: types.Float64Unknown().RefineWithUpperBound(2.1, false),
			min: 2,
		},
		"unknown upper bound exclusive - invalid matches bound": {
			val:         types.Float64Unknown().RefineWithUpperBound(2.1, false),
			min:         2.1,
			expectError: true,
		},
		"unknown upper bound exclusive - invalid greater than bound": {
			val:         types.Float64Unknown().RefineWithUpperBound(2.1, false),
			min:         3,
			expectError: true,
		},
		// Unknown value will be <= 2.1
		"unknown upper bound inclusive - valid less than bound": {
			val: types.Float64Unknown().RefineWithUpperBound(2.1, true),
			min: 2,
		},
		"unknown upper bound inclusive - valid matches bound": {
			val: types.Float64Unknown().RefineWithUpperBound(2.1, true),
			min: 2.1,
		},
		"unknown upper bound inclusive - invalid greater than bound": {
			val:         types.Float64Unknown().RefineWithUpperBound(2.1, true),
			min:         3,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(fmt.Sprintf("ValidateFloat64 - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.Float64Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Float64Response{}
			float64validator.AtLeast(test.min).ValidateFloat64(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterFloat64 - %s", name), func(t *testing.T) {
			t.Parallel()
			request := function.Float64ParameterValidatorRequest{
				Value: test.val,
			}
			response := function.Float64ParameterValidatorResponse{}
			float64validator.AtLeast(test.min).ValidateParameterFloat64(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
