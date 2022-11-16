package listvalidator

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidateList(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request           tfsdk.ValidateAttributeRequest
		expectedListElems []attr.Value
		expectedOk        bool
	}{
		"invalid-type": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.BoolValue(true),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedListElems: nil,
			expectedOk:        false,
		},
		"list-null": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.ListNull(types.StringType),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedListElems: nil,
			expectedOk:        false,
		},
		"list-unknown": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.ListUnknown(types.StringType),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedListElems: nil,
			expectedOk:        false,
		},
		"list-value": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.ListValueMust(
					types.StringType,
					[]attr.Value{
						types.StringValue("first"),
						types.StringValue("second"),
					},
				),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedListElems: []attr.Value{
				types.StringValue("first"),
				types.StringValue("second"),
			},
			expectedOk: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotListElems, gotOk := validateList(context.Background(), testCase.request, &tfsdk.ValidateAttributeResponse{})

			if diff := cmp.Diff(gotListElems, testCase.expectedListElems); diff != "" {
				t.Errorf("unexpected float64 difference: %s", diff)
			}

			if diff := cmp.Diff(gotOk, testCase.expectedOk); diff != "" {
				t.Errorf("unexpected ok difference: %s", diff)
			}
		})
	}
}
