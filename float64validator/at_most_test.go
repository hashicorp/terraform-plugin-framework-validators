package float64validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
)

func TestAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		max         float64
		expectError bool
	}
	tests := map[string]testCase{
		"not a Float64": {
			val:         types.BoolValue(true),
			expectError: true,
		},
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
			request := tfsdk.ValidateAttributeRequest{
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
				AttributeConfig:         test.val,
			}
			response := tfsdk.ValidateAttributeResponse{}
			float64validator.AtMost(test.max).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
