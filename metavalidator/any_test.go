package metavalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestAnyValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val           attr.Value
		anyValidators []tfsdk.AttributeValidator
		expectError   bool
	}
	tests := map[string]testCase{
		"Type mismatch": {
			val: types.Int64{Value: 12},
			anyValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(3),
				stringvalidator.LengthAtLeast(5),
			},
			expectError: true,
		},
		"String invalid": {
			val: types.String{Value: "one"},
			anyValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(4),
				stringvalidator.LengthAtLeast(5),
			},
			expectError: true,
		},
		"String valid": {
			val: types.String{Value: "one"},
			anyValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(5),
				stringvalidator.LengthAtLeast(3),
			},
			expectError: false,
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
			Any(test.anyValidators...).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
