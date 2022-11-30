package float64validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
)

func TestBetweenValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Float64
		min         float64
		max         float64
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Float64": {
			val: types.Float64Unknown(),
			min: 0.90,
			max: 3.10,
		},
		"null Float64": {
			val: types.Float64Null(),
			min: 0.90,
			max: 3.10,
		},
		"valid integer as Float64": {
			val: types.Float64Value(2),
			min: 0.90,
			max: 3.10,
		},
		"valid float as Float64": {
			val: types.Float64Value(2.2),
			min: 0.90,
			max: 3.10,
		},
		"valid float as Float64 min": {
			val: types.Float64Value(0.9),
			min: 0.90,
			max: 3.10,
		},
		"valid float as Float64 max": {
			val: types.Float64Value(3.1),
			min: 0.90,
			max: 3.10,
		},
		"too small float as Float64": {
			val:         types.Float64Value(-1.1111),
			min:         0.90,
			max:         3.10,
			expectError: true,
		},
		"too large float as Float64": {
			val:         types.Float64Value(4.2),
			min:         0.90,
			max:         3.10,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := validator.Float64Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Float64Response{}
			float64validator.Between(test.min, test.max).ValidateFloat64(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
