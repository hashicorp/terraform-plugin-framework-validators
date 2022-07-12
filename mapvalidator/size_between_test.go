package mapvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSizeBetweenValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		min         int
		max         int
		expectError bool
	}
	tests := map[string]testCase{
		"not a Map": {
			val:         types.Bool{Value: true},
			expectError: true,
		},
		"Map unknown": {
			val: types.Map{
				Unknown:  true,
				ElemType: types.StringType,
			},
			expectError: false,
		},
		"Map null": {
			val: types.Map{
				Null:     true,
				ElemType: types.StringType,
			},
			expectError: false,
		},
		"Map size greater than min": {
			val: types.Map{
				ElemType: types.StringType,
				Elems: map[string]attr.Value{
					"one": types.String{Value: "first"},
					"two": types.String{Value: "second"},
				},
			},
			min:         1,
			max:         3,
			expectError: false,
		},
		"Map size equal to min": {
			val: types.Map{
				ElemType: types.StringType,
				Elems: map[string]attr.Value{
					"one": types.String{Value: "first"},
				},
			},
			min:         1,
			max:         3,
			expectError: false,
		},
		"Map size less than max": {
			val: types.Map{
				ElemType: types.StringType,
				Elems: map[string]attr.Value{
					"one": types.String{Value: "first"},
					"two": types.String{Value: "second"},
				},
			},
			min:         1,
			max:         3,
			expectError: false,
		},
		"Map size equal to max": {
			val: types.Map{
				ElemType: types.StringType,
				Elems: map[string]attr.Value{
					"one":   types.String{Value: "first"},
					"two":   types.String{Value: "second"},
					"three": types.String{Value: "third"},
				},
			},
			min:         1,
			max:         3,
			expectError: false,
		},
		"Map size less than min": {
			val: types.Map{
				ElemType: types.StringType,
				Elems:    map[string]attr.Value{},
			},
			min:         1,
			max:         3,
			expectError: true,
		},
		"Map size greater than max": {
			val: types.Map{
				ElemType: types.StringType,
				Elems: map[string]attr.Value{
					"one":   types.String{Value: "first"},
					"two":   types.String{Value: "second"},
					"three": types.String{Value: "third"},
					"four":  types.String{Value: "fourth"},
				},
			},
			min:         1,
			max:         3,
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
			SizeBetween(test.min, test.max).Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
