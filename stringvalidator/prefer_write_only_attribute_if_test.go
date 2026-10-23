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
				ConfigValue:        types.StringValue("oldAttribute value"),
				Path:               path.Root("oldAttribute"),
				PathExpression:     path.MatchRoot("oldAttribute"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"writeOnlyAttribute": schema.StringAttribute{WriteOnly: true},
							"oldAttribute":       schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"writeOnlyAttribute": tftypes.String,
							"oldAttribute":       tftypes.String,
						},
					}, map[string]tftypes.Value{
						"writeOnlyAttribute": tftypes.NewValue(tftypes.String, nil),
						"oldAttribute":       tftypes.NewValue(tftypes.String, "oldAttribute value"),
					}),
				},
			}
			resp := &validator.StringResponse{}
			conditionCalled := false

			stringvalidator.PreferWriteOnlyAttributeIf(
				path.MatchRoot("writeOnlyAttribute"),
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
