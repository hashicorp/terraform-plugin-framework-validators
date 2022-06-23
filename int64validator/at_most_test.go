package int64validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
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
			max: 2,
		},
		"null Int64": {
			val: types.Int64{Null: true},
			max: 2,
		},
		"valid integer as Int64": {
			val: types.Int64{Value: 1},
			max: 2,
		},
		"valid integer as Int64 min": {
			val: types.Int64{Value: 2},
			max: 2,
		},
		"too large integer as Int64": {
			val:         types.Int64{Value: 4},
			max:         2,
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
			int64validator.AtMost(test.max).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
