// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package schemavalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator"
)

func TestPreferWriteOnlyAttribute(t *testing.T) {
	t.Parallel()

	type testCase struct {
		req         schemavalidator.PreferWriteOnlyAttributeRequest
		in          path.Expression
		expWarnings int
		expErrors   int
	}

	testCases := map[string]testCase{
		"base": {
			req: schemavalidator.PreferWriteOnlyAttributeRequest{
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
			},
			in:          path.MatchRoot("writeOnlyAttribute"),
			expWarnings: 1,
		},
		"no-write-only-capability": {
			req: schemavalidator.PreferWriteOnlyAttributeRequest{
				ConfigValue:    types.StringValue("oldAttribute value"),
				Path:           path.Root("oldAttribute"),
				PathExpression: path.MatchRoot("oldAttribute"),
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
			},
			in: path.MatchRoot("writeOnlyAttribute"),
		},
		"old-attribute-is-null": {
			req: schemavalidator.PreferWriteOnlyAttributeRequest{
				ConfigValue:    types.StringNull(),
				Path:           path.Root("oldAttribute"),
				PathExpression: path.MatchRoot("oldAttribute"),
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
						"oldAttribute":       tftypes.NewValue(tftypes.String, nil),
					}),
				},
			},
			in: path.MatchRoot("writeOnlyAttribute"),
		},
		"old-attribute-is-unknown": {
			req: schemavalidator.PreferWriteOnlyAttributeRequest{
				ClientCapabilities: validator.ValidateSchemaClientCapabilities{WriteOnlyAttributesAllowed: true},
				ConfigValue:        types.StringUnknown(),
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
						"oldAttribute":       tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
					}),
				},
			},
			in: path.MatchRoot("writeOnlyAttribute"),
		},
		"matches-no-attribute-in-schema": {
			req: schemavalidator.PreferWriteOnlyAttributeRequest{
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
			},
			in:        path.MatchRoot("writeOnlyAttribute2"),
			expErrors: 1,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := &schemavalidator.PreferWriteOnlyAttributeResponse{}

			schemavalidator.PreferWriteOnlyAttribute{
				WriteOnlyAttribute: test.in,
			}.Validate(context.TODO(), test.req, res)

			if test.expWarnings != res.Diagnostics.WarningsCount() {
				t.Fatalf("expected no warining(s), got %d: %v", res.Diagnostics.WarningsCount(), res.Diagnostics)
			}

			if test.expWarnings > 0 && test.expWarnings != res.Diagnostics.WarningsCount() {
				t.Fatalf("expected %d warning(s), got %d: %v", test.expWarnings, res.Diagnostics.WarningsCount(), res.Diagnostics)
			}

			if test.expErrors == 0 && res.Diagnostics.HasError() {
				t.Fatalf("expected no error(s), got %d: %v", res.Diagnostics.WarningsCount(), res.Diagnostics)
			}

			if test.expErrors > 0 && test.expErrors != res.Diagnostics.ErrorsCount() {
				t.Fatalf("expected %d error(s), got %d: %v", test.expErrors, res.Diagnostics.Errors(), res.Diagnostics)
			}
		})
	}
}
