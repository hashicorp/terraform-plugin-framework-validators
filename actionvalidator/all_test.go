// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actionvalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/actionvalidator"
)

func TestAllValidatorValidateAction(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		validators []action.ConfigValidator
		req        action.ValidateConfigRequest
		expected   *action.ValidateConfigResponse
	}{
		"no-diagnostics": {
			validators: []action.ConfigValidator{
				actionvalidator.ExactlyOneOf(
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				),
				actionvalidator.All(
					actionvalidator.AtLeastOneOf(
						path.MatchRoot("test3"),
						path.MatchRoot("test4"),
					),
					actionvalidator.Conflicting(
						path.MatchRoot("test3"),
						path.MatchRoot("test5"),
					),
				),
			},
			req: action.ValidateConfigRequest{
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
			expected: &action.ValidateConfigResponse{},
		},
		"diagnostics": {
			validators: []action.ConfigValidator{
				actionvalidator.ExactlyOneOf(
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				),
				actionvalidator.All(
					actionvalidator.AtLeastOneOf(
						path.MatchRoot("test3"),
						path.MatchRoot("test4"),
					),
					actionvalidator.Conflicting(
						path.MatchRoot("test3"),
						path.MatchRoot("test5"),
					),
				),
			},
			req: action.ValidateConfigRequest{
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
			expected: &action.ValidateConfigResponse{
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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := &action.ValidateConfigResponse{}

			actionvalidator.Any(testCase.validators...).ValidateAction(context.Background(), testCase.req, got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
