// Copyright (c) HashiCorp, Inc.
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

func TestAllValidatorValidateDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		validators []datasource.ConfigValidator
		req        datasource.ValidateConfigRequest
		expected   *datasource.ValidateConfigResponse
	}{
		"no-diagnostics": {
			validators: []datasource.ConfigValidator{
				datasourcevalidator.ExactlyOneOf(
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				),
				datasourcevalidator.All(
					datasourcevalidator.AtLeastOneOf(
						path.MatchRoot("test3"),
						path.MatchRoot("test4"),
					),
					datasourcevalidator.Conflicting(
						path.MatchRoot("test3"),
						path.MatchRoot("test5"),
					),
				),
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
							"test3": schema.StringAttribute{
								Optional: true,
							},
							"test4": schema.StringAttribute{
								Optional: true,
							},
							"test5": schema.StringAttribute{
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
								"test5": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"test1": tftypes.NewValue(tftypes.String, nil),
							"test2": tftypes.NewValue(tftypes.String, nil),
							"test3": tftypes.NewValue(tftypes.String, "test-value"),
							"test4": tftypes.NewValue(tftypes.String, nil),
							"test5": tftypes.NewValue(tftypes.String, nil),
						},
					),
				},
			},
			expected: &datasource.ValidateConfigResponse{},
		},
		"diagnostics": {
			validators: []datasource.ConfigValidator{
				datasourcevalidator.ExactlyOneOf(
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				),
				datasourcevalidator.All(
					datasourcevalidator.AtLeastOneOf(
						path.MatchRoot("test3"),
						path.MatchRoot("test4"),
					),
					datasourcevalidator.Conflicting(
						path.MatchRoot("test3"),
						path.MatchRoot("test5"),
					),
				),
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
							"test3": schema.StringAttribute{
								Optional: true,
							},
							"test4": schema.StringAttribute{
								Optional: true,
							},
							"test5": schema.StringAttribute{
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
								"test5": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"test1": tftypes.NewValue(tftypes.String, nil),
							"test2": tftypes.NewValue(tftypes.String, nil),
							"test3": tftypes.NewValue(tftypes.String, "test-value"),
							"test4": tftypes.NewValue(tftypes.String, nil),
							"test5": tftypes.NewValue(tftypes.String, "test-value"),
						},
					),
				},
			},
			expected: &datasource.ValidateConfigResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Missing Attribute Configuration",
						"Exactly one of these attributes must be configured: [test1,test2]",
					),
					diag.WithPath(path.Root("test3"),
						diag.NewErrorDiagnostic(
							"Invalid Attribute Combination",
							"These attributes cannot be configured together: [test3,test5]",
						)),
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := &datasource.ValidateConfigResponse{}

			datasourcevalidator.Any(testCase.validators...).ValidateDataSource(context.Background(), testCase.req, got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
