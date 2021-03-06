package setvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSizeAtLeastValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		min         int
		expectError bool
	}
	tests := map[string]testCase{
		"not a Set": {
			val:         types.Bool{Value: true},
			expectError: true,
		},
		"Set unknown": {
			val: types.Set{
				Unknown:  true,
				ElemType: types.StringType,
			},
			expectError: false,
		},
		"Set null": {
			val: types.Set{
				Null:     true,
				ElemType: types.StringType,
			},
			expectError: false,
		},
		"Set size greater than min": {
			val: types.Set{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "first"},
					types.String{Value: "second"},
				},
			},
			min:         1,
			expectError: false,
		},
		"Set size equal to min": {
			val: types.Set{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "first"},
				},
			},
			min:         1,
			expectError: false,
		},
		"Set size less than min": {
			val: types.Set{
				ElemType: types.StringType,
				Elems:    []attr.Value{},
			},
			min:         1,
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
			SizeAtLeast(test.min).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
