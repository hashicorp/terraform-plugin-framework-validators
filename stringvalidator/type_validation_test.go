package stringvalidator

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidateString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request        tfsdk.ValidateAttributeRequest
		expectedString string
		expectedOk     bool
	}{
		"invalid-type": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.BoolValue(true),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedString: "",
			expectedOk:     false,
		},
		"string-null": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.Int64Null(),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedString: "",
			expectedOk:     false,
		},
		"string-value": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.StringValue("test-value"),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedString: "test-value",
			expectedOk:     true,
		},
		"string-unknown": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.Int64Unknown(),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedString: "",
			expectedOk:     false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotInt64, gotOk := validateString(context.Background(), testCase.request, &tfsdk.ValidateAttributeResponse{})

			if diff := cmp.Diff(gotInt64, testCase.expectedString); diff != "" {
				t.Errorf("unexpected float64 difference: %s", diff)
			}

			if diff := cmp.Diff(gotOk, testCase.expectedOk); diff != "" {
				t.Errorf("unexpected ok difference: %s", diff)
			}
		})
	}
}
