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

func TestAnyValidatorValidateListResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		validators []list.ConfigValidator
		req        list.ValidateConfigRequest
		expected   *list.ValidateConfigResponse
	}{
		"no-diagnostics": {
			validators: []list.ConfigValidator{
				listresourcevalidator.ExactlyOneOf(
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				),
				listresourcevalidator.ExactlyOneOf(
					path.MatchRoot("test3"),
					path.MatchRoot("test4"),
				),
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
							"test3": schema.StringAttribute{
								Optional: true,
							},
							"test4": schema.StringAttribute{
								Optional: true,
							},
						},
					},
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test1": tftypes.String,
								"test2": tftypes.String,
								"test3": tftypes.String,
								"test4": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"test1": tftypes.NewValue(tftypes.String, nil),
							"test2": tftypes.NewValue(tftypes.String, nil),
							"test3": tftypes.NewValue(tftypes.String, "test-value"),
							"test4": tftypes.NewValue(tftypes.String, nil),
						},
					),
				},
			},
			expected: &list.ValidateConfigResponse{},
		},
		"diagnostics": {
			validators: []list.ConfigValidator{
				listresourcevalidator.ExactlyOneOf(
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				),
				listresourcevalidator.ExactlyOneOf(
					path.MatchRoot("test3"),
					path.MatchRoot("test4"),
				),
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
							"test3": schema.StringAttribute{
								Optional: true,
							},
							"test4": schema.StringAttribute{
								Optional: true,
							},
						},
					},
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test1": tftypes.String,
								"test2": tftypes.String,
								"test3": tftypes.String,
								"test4": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"test1": tftypes.NewValue(tftypes.String, nil),
							"test2": tftypes.NewValue(tftypes.String, nil),
							"test3": tftypes.NewValue(tftypes.String, nil),
							"test4": tftypes.NewValue(tftypes.String, nil),
						},
					),
				},
			},
			expected: &list.ValidateConfigResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Missing Attribute Configuration",
						"Exactly one of these attributes must be configured: [test1,test2]",
					),
					diag.NewErrorDiagnostic(
						"Missing Attribute Configuration",
						"Exactly one of these attributes must be configured: [test3,test4]",
					),
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := &list.ValidateConfigResponse{}

			listresourcevalidator.Any(testCase.validators...).ValidateListResourceConfig(context.Background(), testCase.req, got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
