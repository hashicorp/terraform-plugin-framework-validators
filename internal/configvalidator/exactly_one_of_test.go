package configvalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestExactlyOneOfValidatorValidate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		validator configvalidator.ExactlyOneOfValidator
		config    tfsdk.Config
		expected  diag.Diagnostics
	}{
		"nil-path-expressions": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: nil,
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test": {
							Optional: true,
							Type:     types.StringType,
						},
					},
				},
				Raw: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					},
					map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Missing Attribute Configuration",
					"Exactly one of these attributes must be configured: []",
				),
			},
		},
		"empty-path-expressions": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test": {
							Optional: true,
							Type:     types.StringType,
						},
					},
				},
				Raw: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					},
					map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Missing Attribute Configuration",
					"Exactly one of these attributes must be configured: []",
				),
			},
		},
		"one-non-existent-path-expression": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("not-test"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test": {
							Optional: true,
							Type:     types.StringType,
						},
					},
				},
				Raw: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					},
					map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Invalid Path Expression for Schema Data",
					"The Terraform Provider unexpectedly matched no paths with the given path expression and current schema data. "+
						"This can happen if the path expression does not correctly follow the schema in structure or types. "+
						"Please report this to the provider developers.\n\n"+
						"Path Expression: not-test",
				),
			},
		},
		"two-non-existent-path-expression": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("not-test1"),
					path.MatchRoot("not-test2"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test": {
							Optional: true,
							Type:     types.StringType,
						},
					},
				},
				Raw: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					},
					map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Invalid Path Expression for Schema Data",
					"The Terraform Provider unexpectedly matched no paths with the given path expression and current schema data. "+
						"This can happen if the path expression does not correctly follow the schema in structure or types. "+
						"Please report this to the provider developers.\n\n"+
						"Path Expression: not-test1",
				),
				diag.NewErrorDiagnostic(
					"Invalid Path Expression for Schema Data",
					"The Terraform Provider unexpectedly matched no paths with the given path expression and current schema data. "+
						"This can happen if the path expression does not correctly follow the schema in structure or types. "+
						"Please report this to the provider developers.\n\n"+
						"Path Expression: not-test2",
				),
			},
		},
		"one-matching-path-expression-null": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test": {
							Optional: true,
							Type:     types.StringType,
						},
						"other": {
							Optional: true,
							Type:     types.StringType,
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
						"test":  tftypes.NewValue(tftypes.String, nil),
						"other": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Missing Attribute Configuration",
					"Exactly one of these attributes must be configured: [test]",
				),
			},
		},
		"one-matching-path-expression-unknown": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test": {
							Optional: true,
							Type:     types.StringType,
						},
						"other": {
							Optional: true,
							Type:     types.StringType,
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
						"test":  tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
						"other": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: nil,
		},
		"one-matching-path-expression-value": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test": {
							Optional: true,
							Type:     types.StringType,
						},
						"other": {
							Optional: true,
							Type:     types.StringType,
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
			expected: nil,
		},
		"two-matching-path-expression-one-null-one-value": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test1": {
							Optional: true,
							Type:     types.StringType,
						},
						"test2": {
							Optional: true,
							Type:     types.StringType,
						},
						"other": {
							Optional: true,
							Type:     types.StringType,
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
						"test2": tftypes.NewValue(tftypes.String, "test-value"),
						"other": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: nil,
		},
		"two-matching-path-expression-one-unknown-one-value": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test1": {
							Optional: true,
							Type:     types.StringType,
						},
						"test2": {
							Optional: true,
							Type:     types.StringType,
						},
						"other": {
							Optional: true,
							Type:     types.StringType,
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
						"test1": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
						"test2": tftypes.NewValue(tftypes.String, "test-value"),
						"other": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: nil,
		},
		"two-matching-path-expression-two-null": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test1": {
							Optional: true,
							Type:     types.StringType,
						},
						"test2": {
							Optional: true,
							Type:     types.StringType,
						},
						"other": {
							Optional: true,
							Type:     types.StringType,
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
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Missing Attribute Configuration",
					"Exactly one of these attributes must be configured: [test1,test2]",
				),
			},
		},
		"two-matching-path-expression-two-unknown": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test1": {
							Optional: true,
							Type:     types.StringType,
						},
						"test2": {
							Optional: true,
							Type:     types.StringType,
						},
						"other": {
							Optional: true,
							Type:     types.StringType,
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
						"test1": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
						"test2": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
						"other": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: nil,
		},
		"two-matching-path-expression-two-value": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test1": {
							Optional: true,
							Type:     types.StringType,
						},
						"test2": {
							Optional: true,
							Type:     types.StringType,
						},
						"other": {
							Optional: true,
							Type:     types.StringType,
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
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test1"),
					"Invalid Attribute Combination",
					"Exactly one of these attributes must be configured: [test1,test2]",
				),
			},
		},
		"three-matching-path-expression-two-value-one-null": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
					path.MatchRoot("test3"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test1": {
							Optional: true,
							Type:     types.StringType,
						},
						"test2": {
							Optional: true,
							Type:     types.StringType,
						},
						"test3": {
							Optional: true,
							Type:     types.StringType,
						},
						"other": {
							Optional: true,
							Type:     types.StringType,
						},
					},
				},
				Raw: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test1": tftypes.String,
							"test2": tftypes.String,
							"test3": tftypes.String,
							"other": tftypes.String,
						},
					},
					map[string]tftypes.Value{
						"test1": tftypes.NewValue(tftypes.String, "test-value"),
						"test2": tftypes.NewValue(tftypes.String, "test-value"),
						"test3": tftypes.NewValue(tftypes.String, nil),
						"other": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test1"),
					"Invalid Attribute Combination",
					"Exactly one of these attributes must be configured: [test1,test2,test3]",
				),
			},
		},
		"three-matching-path-expression-two-value-one-unknown": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
					path.MatchRoot("test3"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test1": {
							Optional: true,
							Type:     types.StringType,
						},
						"test2": {
							Optional: true,
							Type:     types.StringType,
						},
						"test3": {
							Optional: true,
							Type:     types.StringType,
						},
						"other": {
							Optional: true,
							Type:     types.StringType,
						},
					},
				},
				Raw: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test1": tftypes.String,
							"test2": tftypes.String,
							"test3": tftypes.String,
							"other": tftypes.String,
						},
					},
					map[string]tftypes.Value{
						"test1": tftypes.NewValue(tftypes.String, "test-value"),
						"test2": tftypes.NewValue(tftypes.String, "test-value"),
						"test3": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
						"other": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test1"),
					"Invalid Attribute Combination",
					"Exactly one of these attributes must be configured: [test1,test2,test3]",
				),
			},
		},
		"three-matching-path-expression-three-value": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
					path.MatchRoot("test3"),
				},
			},
			config: tfsdk.Config{
				Schema: tfsdk.Schema{
					Attributes: map[string]tfsdk.Attribute{
						"test1": {
							Optional: true,
							Type:     types.StringType,
						},
						"test2": {
							Optional: true,
							Type:     types.StringType,
						},
						"test3": {
							Optional: true,
							Type:     types.StringType,
						},
						"other": {
							Optional: true,
							Type:     types.StringType,
						},
					},
				},
				Raw: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test1": tftypes.String,
							"test2": tftypes.String,
							"test3": tftypes.String,
							"other": tftypes.String,
						},
					},
					map[string]tftypes.Value{
						"test1": tftypes.NewValue(tftypes.String, "test-value"),
						"test2": tftypes.NewValue(tftypes.String, "test-value"),
						"test3": tftypes.NewValue(tftypes.String, "test-value"),
						"other": tftypes.NewValue(tftypes.String, "test-value"),
					},
				),
			},
			expected: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test1"),
					"Invalid Attribute Combination",
					"Exactly one of these attributes must be configured: [test1,test2,test3]",
				),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.validator.Validate(context.Background(), testCase.config)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestExactlyOneOfValidatorValidateDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		validator configvalidator.ExactlyOneOfValidator
		req       datasource.ValidateConfigRequest
		expected  *datasource.ValidateConfigResponse
	}{
		"no-diagnostics": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test"),
				},
			},
			req: datasource.ValidateConfigRequest{
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Optional: true,
								Type:     types.StringType,
							},
							"other": {
								Optional: true,
								Type:     types.StringType,
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
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				},
			},
			req: datasource.ValidateConfigRequest{
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test1": {
								Optional: true,
								Type:     types.StringType,
							},
							"test2": {
								Optional: true,
								Type:     types.StringType,
							},
							"other": {
								Optional: true,
								Type:     types.StringType,
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
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := &datasource.ValidateConfigResponse{}

			testCase.validator.ValidateDataSource(context.Background(), testCase.req, got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestExactlyOneOfValidatorValidateProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		validator configvalidator.ExactlyOneOfValidator
		req       provider.ValidateConfigRequest
		expected  *provider.ValidateConfigResponse
	}{
		"no-diagnostics": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test"),
				},
			},
			req: provider.ValidateConfigRequest{
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Optional: true,
								Type:     types.StringType,
							},
							"other": {
								Optional: true,
								Type:     types.StringType,
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
			expected: &provider.ValidateConfigResponse{},
		},
		"diagnostics": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				},
			},
			req: provider.ValidateConfigRequest{
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test1": {
								Optional: true,
								Type:     types.StringType,
							},
							"test2": {
								Optional: true,
								Type:     types.StringType,
							},
							"other": {
								Optional: true,
								Type:     types.StringType,
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
			expected: &provider.ValidateConfigResponse{
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
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := &provider.ValidateConfigResponse{}

			testCase.validator.ValidateProvider(context.Background(), testCase.req, got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestExactlyOneOfValidatorValidateResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		validator configvalidator.ExactlyOneOfValidator
		req       resource.ValidateConfigRequest
		expected  *resource.ValidateConfigResponse
	}{
		"no-diagnostics": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test"),
				},
			},
			req: resource.ValidateConfigRequest{
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Optional: true,
								Type:     types.StringType,
							},
							"other": {
								Optional: true,
								Type:     types.StringType,
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
			expected: &resource.ValidateConfigResponse{},
		},
		"diagnostics": {
			validator: configvalidator.ExactlyOneOfValidator{
				PathExpressions: path.Expressions{
					path.MatchRoot("test1"),
					path.MatchRoot("test2"),
				},
			},
			req: resource.ValidateConfigRequest{
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test1": {
								Optional: true,
								Type:     types.StringType,
							},
							"test2": {
								Optional: true,
								Type:     types.StringType,
							},
							"other": {
								Optional: true,
								Type:     types.StringType,
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
			expected: &resource.ValidateConfigResponse{
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
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := &resource.ValidateConfigResponse{}

			testCase.validator.ValidateResource(context.Background(), testCase.req, got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
