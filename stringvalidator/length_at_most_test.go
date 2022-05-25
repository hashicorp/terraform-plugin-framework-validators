package stringvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestLengthAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		maxLength   int
		expectError bool
	}
	tests := map[string]testCase{
		"not a String": {
			val:         types.Bool{Value: true},
			expectError: true,
		},
		"unknown String": {
			val:       types.String{Unknown: true},
			maxLength: 1,
		},
		"null String": {
			val:       types.String{Null: true},
			maxLength: 1,
		},
		"valid String": {
			val:       types.String{Value: "ok"},
			maxLength: 2,
		},
		"too long String": {
			val:         types.String{Value: "not ok"},
			maxLength:   5,
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
			LengthAtMost(test.maxLength).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
