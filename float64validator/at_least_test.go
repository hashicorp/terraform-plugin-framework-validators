package float64validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAtLeastValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		min         float64
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
		},
		"null Float64": {
			val: types.Float64{Null: true},
			min: 0.90,
		},
		"valid integer as Float64": {
			val: types.Float64{Value: 2},
			min: 0.90,
		},
		"valid float as Float64": {
			val: types.Float64{Value: 2.2},
			min: 0.90,
		},
		"valid float as Float64 min": {
			val: types.Float64{Value: 0.9},
			min: 0.90,
		},
		"too small float as Float64": {
			val:         types.Float64{Value: -1.1111},
			min:         0.90,
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
			float64validator.AtLeast(test.min).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
