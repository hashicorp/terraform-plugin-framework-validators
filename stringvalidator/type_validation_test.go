package stringvalidator

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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
				AttributeConfig: types.Bool{Value: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedString: "",
			expectedOk:     false,
		},
		"string-null": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Int64{Null: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedString: "",
			expectedOk:     false,
		},
		"string-value": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.String{Value: "test-value"},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedString: "test-value",
			expectedOk:     true,
		},
		"string-unknown": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Int64{Unknown: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
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
