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

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/testvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
)

func TestAnyWithAllWarningsValidatorValidateResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		validators []resource.ConfigValidator
		req        resource.ValidateConfigRequest
		expected   *resource.ValidateConfigResponse
	}{
		"valid": {
			validators: []resource.ConfigValidator{
				resourcevalidator.ExactlyOneOf(
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				),
				resourcevalidator.ExactlyOneOf(
					path.MatchRoot("test3"),
					path.MatchRoot("test4"),
				),
			},
			req: resource.ValidateConfigRequest{
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
			expected: &resource.ValidateConfigResponse{},
		},
		"valid with warning": {
			validators: []resource.ConfigValidator{
				resourcevalidator.All(
					resourcevalidator.ExactlyOneOf(
						path.MatchRoot("test1"),
						path.MatchRoot("test2"),
					),
					testvalidator.WarningResource("failing warning summary", "failing warning details"),
				),
				resourcevalidator.All(
					resourcevalidator.ExactlyOneOf(
						path.MatchRoot("test3"),
						path.MatchRoot("test4"),
					),
					testvalidator.WarningResource("passing warning summary", "passing warning details"),
				),
			},
			req: resource.ValidateConfigRequest{
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
			expected: &resource.ValidateConfigResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewWarningDiagnostic("failing warning summary", "failing warning details"),
					diag.NewWarningDiagnostic("passing warning summary", "passing warning details"),
				},
			},
		},
		"invalid": {
			validators: []resource.ConfigValidator{
				resourcevalidator.ExactlyOneOf(
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				),
				resourcevalidator.ExactlyOneOf(
					path.MatchRoot("test3"),
					path.MatchRoot("test4"),
				),
			},
			req: resource.ValidateConfigRequest{
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
			expected: &resource.ValidateConfigResponse{
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
