package stringvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestLengthBetweenValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		minLength   int
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
			minLength: 1,
			maxLength: 3,
		},
		"null String": {
			val:       types.String{Null: true},
			minLength: 1,
			maxLength: 3,
		},
		"valid String": {
			val:       types.String{Value: "ok"},
			minLength: 1,
			maxLength: 3,
		},
		"too long String": {
			val:         types.String{Value: "not ok"},
			minLength:   1,
			maxLength:   3,
			expectError: true,
		},
		"too short String": {
			val:         types.String{Value: ""},
			minLength:   1,
			maxLength:   3,
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
			LengthBetween(test.minLength, test.maxLength).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
