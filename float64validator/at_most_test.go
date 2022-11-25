package float64validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
)

func TestAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Float64
		max         float64
		expectError bool
	}
	tests := map[string]testCase{
		"unknown Float64": {
			val: types.Float64Unknown(),
			max: 2.00,
		},
		"null Float64": {
			val: types.Float64Null(),
			max: 2.00,
		},
		"valid integer as Float64": {
			val: types.Float64Value(1),
			max: 2.00,
		},
		"valid float as Float64": {
			val: types.Float64Value(1.1),
			max: 2.00,
		},
		"valid float as Float64 max": {
			val: types.Float64Value(2.0),
			max: 2.00,
		},
		"too large float as Float64": {
			val:         types.Float64Value(3.0),
			max:         2.00,
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
			float64validator.AtMost(test.max).ValidateFloat64(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
