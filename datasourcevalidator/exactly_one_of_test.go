// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package datasourcevalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
)

func TestExactlyOneOf(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		pathExpressions path.Expressions
		req             datasource.ValidateConfigRequest
		expected        *datasource.ValidateConfigResponse
	}{
		"no-diagnostics": {
			pathExpressions: path.Expressions{
				path.MatchRoot("test"),
			},
			req: datasource.ValidateConfigRequest{
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"test": schema.StringAttribute{
								Optional: true,
							},
							"other": schema.StringAttribute{
								Optional: true,
							},
						},
					},
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test":  tftypes.String,
								"other": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"test":  tftypes.NewValue(tftypes.String, "test-value"),
							"other": tftypes.NewValue(tftypes.String, "test-value"),
						},
					),
				},
			},
			expected: &datasource.ValidateConfigResponse{},
		},
		"diagnostics": {
			pathExpressions: path.Expressions{
				path.MatchRoot("test1"),
				path.MatchRoot("test2"),
			},
			req: datasource.ValidateConfigRequest{
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"test1": schema.StringAttribute{
								Optional: true,
							},
							"test2": schema.StringAttribute{
								Optional: true,
							},
							"other": schema.StringAttribute{
								Optional: true,
							},
						},
					},
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test1": tftypes.String,
								"test2": tftypes.String,
								"other": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"test1": tftypes.NewValue(tftypes.String, "test-value"),
							"test2": tftypes.NewValue(tftypes.String, "test-value"),
							"other": tftypes.NewValue(tftypes.String, "test-value"),
						},
					),
				},
			},
			expected: &datasource.ValidateConfigResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						path.Root("test1"),
						"Invalid Attribute Combination",
						"Exactly one of these attributes must be configured: [test1,test2]",
					),
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			validator := datasourcevalidator.ExactlyOneOf(testCase.pathExpressions...)
			got := &datasource.ValidateConfigResponse{}

			validator.ValidateDataSource(context.Background(), testCase.req, got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
