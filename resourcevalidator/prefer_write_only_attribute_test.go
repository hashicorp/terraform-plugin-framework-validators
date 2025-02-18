// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcevalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
)

func TestPreferWriteOnlyAttribute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		validators []resource.ConfigValidator
		req        resource.ValidateConfigRequest
		expected   *resource.ValidateConfigResponse
	}{
		"valid-warning-diag": {
			validators: []resource.ConfigValidator{
				resourcevalidator.PreferWriteOnlyAttribute(
					path.MatchRoot("oldAttribute1"),
					path.MatchRoot("writeOnlyAttribute1"),
				),
				resourcevalidator.PreferWriteOnlyAttribute(
					path.MatchRoot("oldAttribute2"),
					path.MatchRoot("writeOnlyAttribute2"),
				),
			},
			req: resource.ValidateConfigRequest{
				ClientCapabilities: resource.ValidateConfigClientCapabilities{WriteOnlyAttributesAllowed: true},
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"oldAttribute1": schema.StringAttribute{
								Optional: true,
							},
							"writeOnlyAttribute1": schema.StringAttribute{
								Optional:  true,
								WriteOnly: true,
							},
							"oldAttribute2": schema.StringAttribute{
								Optional: true,
							},
							"writeOnlyAttribute2": schema.StringAttribute{
								Optional:  true,
								WriteOnly: true,
							},
						},
					},
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"oldAttribute1":       tftypes.String,
								"writeOnlyAttribute1": tftypes.String,
								"oldAttribute2":       tftypes.String,
								"writeOnlyAttribute2": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"oldAttribute1":       tftypes.NewValue(tftypes.String, nil),
							"writeOnlyAttribute1": tftypes.NewValue(tftypes.String, nil),
							"oldAttribute2":       tftypes.NewValue(tftypes.String, "test-value"),
							"writeOnlyAttribute2": tftypes.NewValue(tftypes.String, nil),
						},
					),
				},
			},
			expected: &resource.ValidateConfigResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeWarningDiagnostic(path.Root("oldAttribute2"),
						"Available Write-Only Attribute Alternative",
						"The attribute has a WriteOnly version writeOnlyAttribute2 available. "+
							"Use the WriteOnly version of the attribute when possible."),
				},
			},
		},
		"valid-no-client-capabilities": {
			validators: []resource.ConfigValidator{
				resourcevalidator.PreferWriteOnlyAttribute(
					path.MatchRoot("oldAttribute1"),
					path.MatchRoot("writeOnlyAttribute1"),
				),
				resourcevalidator.PreferWriteOnlyAttribute(
					path.MatchRoot("oldAttribute2"),
					path.MatchRoot("writeOnlyAttribute2"),
				),
			},
			req: resource.ValidateConfigRequest{
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"oldAttribute1": schema.StringAttribute{
								Optional: true,
							},
							"writeOnlyAttribute1": schema.StringAttribute{
								Optional: true,
							},
							"oldAttribute2": schema.StringAttribute{
								Optional: true,
							},
							"writeOnlyAttribute2": schema.StringAttribute{
								Optional: true,
							},
						},
					},
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"oldAttribute1":       tftypes.String,
								"writeOnlyAttribute1": tftypes.String,
								"oldAttribute2":       tftypes.String,
								"writeOnlyAttribute2": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"oldAttribute1":       tftypes.NewValue(tftypes.String, nil),
							"writeOnlyAttribute1": tftypes.NewValue(tftypes.String, nil),
							"oldAttribute2":       tftypes.NewValue(tftypes.String, "test-value"),
							"writeOnlyAttribute2": tftypes.NewValue(tftypes.String, nil),
						},
					),
				},
			},
			expected: &resource.ValidateConfigResponse{},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := &resource.ValidateConfigResponse{}

			resourcevalidator.AnyWithAllWarnings(testCase.validators...).ValidateResource(context.Background(), testCase.req, got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
