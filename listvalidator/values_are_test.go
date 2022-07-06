package listvalidator

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
		"not List": {
			val: types.Set{
				ElemType: types.StringType,
			},
			expectError: true,
		},
		"List unknown": {
			val: types.List{
				Unknown:  true,
				ElemType: types.StringType,
			},
			expectError: false,
		},
		"List null": {
			val: types.List{
				Null:     true,
				ElemType: types.StringType,
			},
			expectError: false,
		},
		"List elems invalid": {
			val: types.List{
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
		"List elems invalid for second validator": {
			val: types.List{
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
		"List elems wrong type for validator": {
			val: types.List{
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
		"List elems valid": {
			val: types.List{
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
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
				AttributeConfig:         test.val,
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
