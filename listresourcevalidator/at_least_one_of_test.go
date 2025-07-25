// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listresourcevalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-validators/listresourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAtLeastOneOf(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		pathExpressions path.Expressions
		req             list.ValidateConfigRequest
		expected        *list.ValidateConfigResponse
	}{
		"no-diagnostics": {
			pathExpressions: path.Expressions{
				path.MatchRoot("test"),
			},
			req: list.ValidateConfigRequest{
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
			expected: &list.ValidateConfigResponse{},
		},
		"diagnostics": {
			pathExpressions: path.Expressions{
				path.MatchRoot("test1"),
				path.MatchRoot("test2"),
			},
			req: list.ValidateConfigRequest{
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
							"test1": tftypes.NewValue(tftypes.String, nil),
							"test2": tftypes.NewValue(tftypes.String, nil),
							"other": tftypes.NewValue(tftypes.String, "test-value"),
						},
					),
				},
			},
			expected: &list.ValidateConfigResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Missing Attribute Configuration",
						"At least one of these attributes must be configured: [test1,test2]",
					),
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			validator := listresourcevalidator.AtLeastOneOf(testCase.pathExpressions...)
			got := &list.ValidateConfigResponse{}

			validator.ValidateListResourceConfig(context.Background(), testCase.req, got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
