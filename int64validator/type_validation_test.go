package int64validator

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidateInt(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request             tfsdk.ValidateAttributeRequest
		expectedInt64       int64
		expectedOk          bool
		expectedDiagSummary string
		expectedDiagDetail  string
	}{
		"invalid-type": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.BoolValue(true),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedInt64:       0.0,
			expectedOk:          false,
			expectedDiagSummary: "Value Conversion Error",
			expectedDiagDetail:  "An unexpected error was encountered trying to convert into a Terraform value. This is always an error in the provider. Please report the following to the provider developer:\n\nCannot use attr.Value types.Int64, only types.Bool is supported because types.primitive is the type in the schema",
		},
		"int64-null": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.Int64Null(),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedInt64:       0.0,
			expectedOk:          false,
			expectedDiagSummary: "",
			expectedDiagDetail:  "",
		},
		"int64-value": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.Int64Value(123),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedInt64: 123,
			expectedOk:    true,
		},
		"int64-unknown": {
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.Int64Unknown(),
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
			},
			expectedInt64:       0.0,
			expectedOk:          false,
			expectedDiagSummary: "",
			expectedDiagDetail:  "",
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res := tfsdk.ValidateAttributeResponse{}
			gotInt64, gotOk := validateInt(context.Background(), testCase.request, &res)

			if res.Diagnostics.HasError() {
				if res.Diagnostics.ErrorsCount() != 1 {
					t.Errorf("expected an error but found none")
				} else {
					if diff := cmp.Diff(res.Diagnostics[0].Summary(), testCase.expectedDiagSummary); diff != "" {
						t.Errorf("unexpected diagnostic summary difference: %s", diff)
					}

					if diff := cmp.Diff(res.Diagnostics[0].Detail(), testCase.expectedDiagDetail); diff != "" {
						t.Errorf("unexpected diagnostic summary difference: %s", diff)
					}
				}
			}

			if diff := cmp.Diff(gotInt64, testCase.expectedInt64); diff != "" {
				t.Errorf("unexpected int64 difference: %s", diff)
			}

			if diff := cmp.Diff(gotOk, testCase.expectedOk); diff != "" {
				t.Errorf("unexpected ok difference: %s", diff)
			}
		})
	}
}
