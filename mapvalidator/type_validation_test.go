package mapvalidator

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidateMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request     tfsdk.ValidateAttributeRequest
		expectedMap map[string]attr.Value
		expectedOk  bool
	}{
		"invalid-type": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.Bool{Value: true},
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedMap: nil,
			expectedOk:  false,
		},
		"map-null": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.Map{Null: true},
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedMap: nil,
			expectedOk:  false,
		},
		"map-unknown": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.Map{Unknown: true},
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedMap: nil,
			expectedOk:  false,
		},
		"map-value": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Map{
					ElemType: types.StringType,
					Elems: map[string]attr.Value{
						"one": types.String{Value: "first"},
						"two": types.String{Value: "second"},
					},
				},
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedMap: map[string]attr.Value{
				"one": types.String{Value: "first"},
				"two": types.String{Value: "second"},
			},
			expectedOk: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotMapElems, gotOk := validateMap(context.Background(), testCase.request, &tfsdk.ValidateAttributeResponse{})

			if diff := cmp.Diff(gotMapElems, testCase.expectedMap); diff != "" {
				t.Errorf("unexpected map difference: %s", diff)
			}

			if diff := cmp.Diff(gotOk, testCase.expectedOk); diff != "" {
				t.Errorf("unexpected ok difference: %s", diff)
			}
		})
	}
}
