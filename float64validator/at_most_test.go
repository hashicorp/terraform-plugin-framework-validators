package float64validator

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		max         float64
		expectError bool
	}
	tests := map[string]testCase{
		"not a number": {
			val:         types.Bool{Value: true},
			max:         2.00,
			expectError: true,
		},
		"unknown number": {
			val:         types.Number{Unknown: true},
			max:         2.00,
			expectError: true,
		},
		"null number": {
			val:         types.Float64{Null: true},
			max:         2.00,
			expectError: true,
		},
		"valid integer as Number": {
			val: types.Number{Value: big.NewFloat(1)},
			max: 2.00,
		},
		"valid integer as Float64": {
			val: types.Float64{Value: 1},
			max: 2.00,
		},
		"valid float as Number": {
			val: types.Number{Value: big.NewFloat(1.1)},
			max: 2.00,
		},
		"valid float as Float64": {
			val: types.Float64{Value: 1.1},
			max: 2.00,
		},
		"valid float as Number max": {
			val: types.Number{Value: big.NewFloat(2.0)},
			max: 2.00,
		},
		"valid float as Float64 max": {
			val: types.Float64{Value: 2.0},
			max: 2.00,
		},
		"too large float as Number": {
			val:         types.Number{Value: big.NewFloat(3.0)},
			max:         2.00,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := tfsdk.ValidateAttributeRequest{
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
				AttributeConfig: test.val,
			}
			response := tfsdk.ValidateAttributeResponse{}
			AtMost(test.max).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
