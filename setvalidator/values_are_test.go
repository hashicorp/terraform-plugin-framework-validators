package setvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestValuesAreValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val                 attr.Value
		valuesAreValidators []tfsdk.AttributeValidator
		expectError         bool
	}
	tests := map[string]testCase{
		"not Set": {
			val: types.Map{
				ElemType: types.StringType,
			},
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
		"Set elems invalid": {
			val: types.Set{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "first"},
					types.String{Value: "second"},
				},
			},
			valuesAreValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(6),
			},
			expectError: true,
		},
		"Set elems invalid for second validator": {
			val: types.Set{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "first"},
					types.String{Value: "second"},
				},
			},
			valuesAreValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(2),
				stringvalidator.LengthAtLeast(6),
			},
			expectError: true,
		},
		"Set elems wrong type for validator": {
			val: types.Set{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "first"},
					types.String{Value: "second"},
				},
			},
			valuesAreValidators: []tfsdk.AttributeValidator{
				int64validator.AtLeast(6),
			},
			expectError: true,
		},
		"Set elems valid": {
			val: types.Set{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "first"},
					types.String{Value: "second"},
				},
			},
			valuesAreValidators: []tfsdk.AttributeValidator{
				stringvalidator.LengthAtLeast(5),
			},
			expectError: false,
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
			ValuesAre(test.valuesAreValidators...).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
