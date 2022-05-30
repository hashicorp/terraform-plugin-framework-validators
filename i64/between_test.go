package i64

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
		min         int64
		max         int64
		expectError bool
	}
	tests := map[string]testCase{
		"not an Int64": {
			val:         types.Bool{Value: true},
			expectError: true,
		},
		"unknown Int64": {
			val: types.Int64{Unknown: true},
			min: 1,
			max: 3,
		},
		"null Int64": {
			val: types.Int64{Null: true},
			min: 1,
			max: 3,
		},
		"valid integer as Int64": {
			val: types.Int64{Value: 2},
			min: 1,
			max: 3,
		},
		"valid integer as Int64 min": {
			val: types.Int64{Value: 1},
			min: 1,
			max: 3,
		},
		"valid integer as Int64 max": {
			val: types.Int64{Value: 3},
			min: 1,
			max: 3,
		},
		"too small integer as Int64": {
			val:         types.Int64{Value: -1},
			min:         1,
			max:         3,
			expectError: true,
		},
		"too large integer as Int64": {
			val:         types.Int64{Value: 42},
			min:         1,
			max:         3,
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
