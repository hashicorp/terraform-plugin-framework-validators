// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float32validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float32validator"
)

func TestBetweenValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Float32
		min         float32
		max         float32
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Float32": {
			val: types.Float32Unknown(),
			min: 0.90,
			max: 3.10,
		},
		"null Float32": {
			val: types.Float32Null(),
			min: 0.90,
			max: 3.10,
		},
		"valid integer as Float32": {
			val: types.Float32Value(2),
			min: 0.90,
			max: 3.10,
		},
		"valid float as Float32": {
			val: types.Float32Value(2.2),
			min: 0.90,
			max: 3.10,
		},
		"valid float as Float32 min": {
			val: types.Float32Value(0.9),
			min: 0.90,
			max: 3.10,
		},
		"valid float as Float32 max": {
			val: types.Float32Value(3.1),
			min: 0.90,
			max: 3.10,
		},
		"too small float as Float32": {
			val:         types.Float32Value(-1.1111),
			min:         0.90,
			max:         3.10,
			expectError: true,
		},
		"too large float as Float32": {
			val:         types.Float32Value(4.2),
			min:         0.90,
			max:         3.10,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.Float32Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Float32Response{}
			float32validator.Between(test.min, test.max).ValidateFloat32(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
