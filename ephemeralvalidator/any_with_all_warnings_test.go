// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package ephemeralvalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/ephemeralvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/testvalidator"
)

func TestAnyWithAllWarningsValidatorValidateEphemeralResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		validators []ephemeral.ConfigValidator
		req        ephemeral.ValidateConfigRequest
		expected   *ephemeral.ValidateConfigResponse
	}{
		"valid": {
			validators: []ephemeral.ConfigValidator{
				ephemeralvalidator.ExactlyOneOf(
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				),
				ephemeralvalidator.ExactlyOneOf(
					path.MatchRoot("test3"),
					path.MatchRoot("test4"),
				),
			},
			req: ephemeral.ValidateConfigRequest{
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
			expected: &ephemeral.ValidateConfigResponse{},
		},
		"valid with warning": {
			validators: []ephemeral.ConfigValidator{
				ephemeralvalidator.All(
					ephemeralvalidator.ExactlyOneOf(
						path.MatchRoot("test1"),
						path.MatchRoot("test2"),
					),
					testvalidator.WarningEphemeralResource("failing warning summary", "failing warning details"),
				),
				ephemeralvalidator.All(
					ephemeralvalidator.ExactlyOneOf(
						path.MatchRoot("test3"),
						path.MatchRoot("test4"),
					),
					testvalidator.WarningEphemeralResource("passing warning summary", "passing warning details"),
				),
			},
			req: ephemeral.ValidateConfigRequest{
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
			expected: &ephemeral.ValidateConfigResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewWarningDiagnostic("failing warning summary", "failing warning details"),
					diag.NewWarningDiagnostic("passing warning summary", "passing warning details"),
				},
			},
		},
		"invalid": {
			validators: []ephemeral.ConfigValidator{
				ephemeralvalidator.ExactlyOneOf(
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				),
				ephemeralvalidator.ExactlyOneOf(
					path.MatchRoot("test3"),
					path.MatchRoot("test4"),
				),
			},
			req: ephemeral.ValidateConfigRequest{
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
			expected: &ephemeral.ValidateConfigResponse{
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

			got := &ephemeral.ValidateConfigResponse{}

			ephemeralvalidator.AnyWithAllWarnings(testCase.validators...).ValidateEphemeralResource(context.Background(), testCase.req, got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
