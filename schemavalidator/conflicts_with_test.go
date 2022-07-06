package schemavalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestConflictsWithValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		req       tfsdk.ValidateAttributeRequest
		in        path.Expressions
		expErrors int
	}

	testCases := map[string]testCase{
		"base": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.String{Value: "bar value"},
				AttributePath:           path.Root("bar"),
				AttributePathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"foo": {
								Type: types.Int64Type,
							},
							"bar": {
								Type: types.StringType,
							},
							"baz": {
								Type: types.Int64Type,
							},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.String,
							"baz": tftypes.Number,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, 42),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
						"baz": tftypes.NewValue(tftypes.Number, 43),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
				path.MatchRoot("baz"),
			},
			expErrors: 2,
		},
		"conflicting-is-nil": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.String{Value: "bar value"},
				AttributePath:           path.Root("bar"),
				AttributePathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"foo": {
								Type: types.Int64Type,
							},
							"bar": {
								Type: types.StringType,
							},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, nil),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
		},
		"error_conflicting-is-unknown": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.String{Value: "bar value"},
				AttributePath:           path.Root("bar"),
				AttributePathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"foo": {
								Type: types.Int64Type,
							},
							"bar": {
								Type: types.StringType,
							},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
			expErrors: 1,
		},
		"self-is-null": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.String{Null: true},
				AttributePath:           path.Root("bar"),
				AttributePathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"foo": {
								Type: types.Int64Type,
							},
							"bar": {
								Type: types.StringType,
							},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, 42),
						"bar": tftypes.NewValue(tftypes.String, nil),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
		},
		"error_allow-duplicate-input": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.String{Value: "bar value"},
				AttributePath:           path.Root("bar"),
				AttributePathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"foo": {
								Type: types.Int64Type,
							},
							"bar": {
								Type: types.StringType,
							},
							"baz": {
								Type: types.Int64Type,
							},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.String,
							"baz": tftypes.Number,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, 42),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
						"baz": tftypes.NewValue(tftypes.Number, 43),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
				path.MatchRoot("bar"),
				path.MatchRoot("baz"),
			},
			expErrors: 2,
		},
		"error_unknowns": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.String{Value: "bar value"},
				AttributePath:           path.Root("bar"),
				AttributePathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"foo": {
								Type: types.Int64Type,
							},
							"bar": {
								Type: types.StringType,
							},
							"baz": {
								Type: types.Int64Type,
							},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.String,
							"baz": tftypes.Number,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
						"baz": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
				path.MatchRoot("baz"),
			},
			expErrors: 2,
		},
		"matches-no-attribute-in-schema": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig:         types.String{Value: "bar value"},
				AttributePath:           path.Root("bar"),
				AttributePathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"foo": {
								Type: types.Int64Type,
							},
							"bar": {
								Type: types.StringType,
							},
							"baz": {
								Type: types.Int64Type,
							},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.String,
							"baz": tftypes.Number,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, 42),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
						"baz": tftypes.NewValue(tftypes.Number, 43),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("fooz"),
				path.MatchRoot("barz"),
			},
			expErrors: 2,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			res := tfsdk.ValidateAttributeResponse{}

			schemavalidator.ConflictsWith(test.in...).Validate(context.TODO(), test.req, &res)

			if test.expErrors > 0 && !res.Diagnostics.HasError() {
				t.Fatal("expected error(s), got none")
			}

			if test.expErrors > 0 && test.expErrors != res.Diagnostics.ErrorsCount() {
				t.Fatalf("expected %d error(s), got %d: %v", test.expErrors, res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}

			if test.expErrors == 0 && res.Diagnostics.HasError() {
				t.Fatalf("expected no error(s), got %d: %v", res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}
		})
	}
}
