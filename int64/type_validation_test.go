package int64

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestValidateInt(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request       tfsdk.ValidateAttributeRequest
		expectedInt64 int64
		expectedOk    bool
	}{
		"invalid-type": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Bool{Value: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedInt64: 0.0,
			expectedOk:    false,
		},
		"int64-null": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Int64{Null: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedInt64: 0.0,
			expectedOk:    false,
		},
		"int64-value": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Int64{Value: 123},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedInt64: 123,
			expectedOk:    true,
		},
		"int64-unknown": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Int64{Unknown: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedInt64: 0.0,
			expectedOk:    false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotInt64, gotOk := validateInt(context.Background(), testCase.request, &tfsdk.ValidateAttributeResponse{})

			if diff := cmp.Diff(gotInt64, testCase.expectedInt64); diff != "" {
				t.Errorf("unexpected float64 difference: %s", diff)
			}

			if diff := cmp.Diff(gotOk, testCase.expectedOk); diff != "" {
				t.Errorf("unexpected ok difference: %s", diff)
			}
		})
	}
}
