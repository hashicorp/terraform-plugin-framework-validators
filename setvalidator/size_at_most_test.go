package setvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSizeAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		max         int
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
		"Set size less than max": {
			val: types.Set{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "first"},
				},
			},
			max:         2,
			expectError: false,
		},
		"Set size equal to max": {
			val: types.Set{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "first"},
					types.String{Value: "second"},
				},
			},
			max:         2,
			expectError: false,
		},
		"Set size greater than max": {
			val: types.Set{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "first"},
					types.String{Value: "second"},
					types.String{Value: "third"},
				}},
			max:         2,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := tfsdk.ValidateAttributeRequest{
				AttributePath:   path.Root("test"),
				AttributeConfig: test.val,
			}
			response := tfsdk.ValidateAttributeResponse{}
			SizeAtMost(test.max).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
