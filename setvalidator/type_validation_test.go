package setvalidator

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestValidateSet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request          tfsdk.ValidateAttributeRequest
		expectedSetElems []attr.Value
		expectedOk       bool
	}{
		"invalid-type": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Bool{Value: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedSetElems: nil,
			expectedOk:       false,
		},
		"set-null": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Set{Null: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedSetElems: nil,
			expectedOk:       false,
		},
		"set-unknown": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Set{Unknown: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedSetElems: nil,
			expectedOk:       false,
		},
		"set-value": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Set{
					ElemType: types.StringType,
					Elems: []attr.Value{
						types.String{Value: "first"},
						types.String{Value: "second"},
					},
				},
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedSetElems: []attr.Value{
				types.String{Value: "first"},
				types.String{Value: "second"},
			},
			expectedOk: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotSetElems, gotOk := validateSet(context.Background(), testCase.request, &tfsdk.ValidateAttributeResponse{})

			if diff := cmp.Diff(gotSetElems, testCase.expectedSetElems); diff != "" {
				t.Errorf("unexpected set difference: %s", diff)
			}

			if diff := cmp.Diff(gotOk, testCase.expectedOk); diff != "" {
				t.Errorf("unexpected ok difference: %s", diff)
			}
		})
	}
}
