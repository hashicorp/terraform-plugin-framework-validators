package float64validator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestBetweenValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		min         float64
		max         float64
		expectError bool
	}
	tests := map[string]testCase{
		"not a Float64": {
			val:         types.Bool{Value: true},
			expectError: true,
		},
		"unknown Float64": {
			val: types.Float64{Unknown: true},
			min: 0.90,
			max: 3.10,
		},
		"null Float64": {
			val: types.Float64{Null: true},
			min: 0.90,
			max: 3.10,
		},
		"valid integer as Float64": {
			val: types.Float64{Value: 2},
			min: 0.90,
			max: 3.10,
		},
		"valid float as Float64": {
			val: types.Float64{Value: 2.2},
			min: 0.90,
			max: 3.10,
		},
		"valid float as Float64 min": {
			val: types.Float64{Value: 0.9},
			min: 0.90,
			max: 3.10,
		},
		"valid float as Float64 max": {
			val: types.Float64{Value: 3.1},
			min: 0.90,
			max: 3.10,
		},
		"too small float as Float64": {
			val:         types.Float64{Value: -1.1111},
			min:         0.90,
			max:         3.10,
			expectError: true,
		},
		"too large float as Float64": {
			val:         types.Float64{Value: 4.2},
			min:         0.90,
			max:         3.10,
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
			Between(test.min, test.max).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
