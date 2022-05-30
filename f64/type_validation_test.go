package f64

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestValidateFloat(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request         tfsdk.ValidateAttributeRequest
		expectedFloat64 float64
		expectedOk      bool
	}{
		"invalid-type": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Bool{Value: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedFloat64: 0.0,
			expectedOk:      false,
		},
		"float64-null": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Float64{Null: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedFloat64: 0.0,
			expectedOk:      false,
		},
		"float64-value": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Float64{Value: 1.2},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedFloat64: 1.2,
			expectedOk:      true,
		},
		"float64-unknown": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Float64{Unknown: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("test"),
			},
			expectedFloat64: 0.0,
			expectedOk:      false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotFloat64, gotOk := validateFloat(context.Background(), testCase.request, &tfsdk.ValidateAttributeResponse{})

			if diff := cmp.Diff(gotFloat64, testCase.expectedFloat64); diff != "" {
				t.Errorf("unexpected float64 difference: %s", diff)
			}

			if diff := cmp.Diff(gotOk, testCase.expectedOk); diff != "" {
				t.Errorf("unexpected ok difference: %s", diff)
			}
		})
	}
}
