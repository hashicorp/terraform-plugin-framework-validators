// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package stringvalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestPreferWriteOnlyAttributeIf(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		condition   bool
		ifFuncError bool
		expWarnings int
		expErrors   int
	}{
		"condition-true": {
			condition:   true,
			expWarnings: 1,
		},
		"condition-false": {},
		"condition-diagnostic": {
			ifFuncError: true,
			expErrors:   1,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{
				ClientCapabilities: validator.ValidateSchemaClientCapabilities{WriteOnlyAttributesAllowed: true},
				ConfigValue:        types.StringValue("configured value"),
				Path:               path.Root("value"),
				PathExpression:     path.MatchRoot("value"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"value_wo": schema.StringAttribute{WriteOnly: true},
							"value":    schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"value_wo": tftypes.String,
							"value":    tftypes.String,
						},
					}, map[string]tftypes.Value{
						"value_wo": tftypes.NewValue(tftypes.String, nil),
						"value":    tftypes.NewValue(tftypes.String, "configured value"),
					}),
				},
			}
			resp := &validator.StringResponse{}
			conditionCalled := false

			stringvalidator.PreferWriteOnlyAttributeIf(
				path.MatchRoot("value_wo"),
				func(ctx context.Context, gotReq validator.StringRequest, ifResp *stringvalidator.PreferWriteOnlyAttributeIfFuncResponse) {
					conditionCalled = true
					ifResp.Warn = testCase.condition
					if testCase.ifFuncError {
						ifResp.Diagnostics.Append(diag.NewErrorDiagnostic("Conditional Error", "error"))
					}
				},
			).ValidateString(context.Background(), req, resp)

			if !conditionCalled {
				t.Fatal("expected condition function to be called")
			}

			if got := resp.Diagnostics.WarningsCount(); got != testCase.expWarnings {
				t.Fatalf("expected %d warning(s), got %d: %v", testCase.expWarnings, got, resp.Diagnostics)
			}

			if got := resp.Diagnostics.ErrorsCount(); got != testCase.expErrors {
				t.Fatalf("expected %d error(s), got %d: %v", testCase.expErrors, got, resp.Diagnostics)
			}
		})
	}
}
