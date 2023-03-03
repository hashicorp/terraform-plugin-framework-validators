package schemavalidator_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator"
)

func TestNotEqualValidatorValidate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		req       schemavalidator.NotEqualValidatorRequest
		in        path.Expressions
		expErrors int
	}

	testCases := map[string]testCase{
		"base": {
			req: schemavalidator.NotEqualValidatorRequest{
				ConfigValue:    types.StringValue("bar value"),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, "foo value"),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
		},
		"self-is-null": {
			req: schemavalidator.NotEqualValidatorRequest{
				ConfigValue:    types.StringNull(),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, "foo value"),
						"bar": tftypes.NewValue(tftypes.String, nil),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
		},
		"self-is-unknown": {
			req: schemavalidator.NotEqualValidatorRequest{
				ConfigValue:    types.StringUnknown(),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, "foo value"),
						"bar": tftypes.NewValue(tftypes.String, nil),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
		},
		"other-is-null": {
			req: schemavalidator.NotEqualValidatorRequest{
				ConfigValue:    types.StringValue("bar value"),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, nil),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
		},
		"compatible-types": {
			req: schemavalidator.NotEqualValidatorRequest{
				ConfigValue:    types.Int64Value(42),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.NumberAttribute{},
							"bar": schema.Int64Attribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.Number,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, 43),
						"bar": tftypes.NewValue(tftypes.Number, 42),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
		},
		"incompatible-types": {
			req: schemavalidator.NotEqualValidatorRequest{
				ConfigValue:    types.Int64Value(42),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.Int64Attribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.Number,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, "foo value"),
						"bar": tftypes.NewValue(tftypes.Number, 42),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
		},
		"error_same-type-same-value": {
			req: schemavalidator.NotEqualValidatorRequest{
				ConfigValue:    types.StringValue("bar value"),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, "bar value"),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
			expErrors: 1,
		},
		"error_multiple-same-type-same-value": {
			req: schemavalidator.NotEqualValidatorRequest{
				ConfigValue:    types.StringValue("bar value"),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
							"baz": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
							"baz": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, "bar value"),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
						"baz": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
				path.MatchRoot("baz"),
			},
			expErrors: 2,
		},
		"error_compatible-types-same-value": {
			req: schemavalidator.NotEqualValidatorRequest{
				ConfigValue:    types.NumberValue(big.NewFloat(float64(42))),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.Int64Attribute{},
							"bar": schema.NumberAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.Number,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, 42),
						"bar": tftypes.NewValue(tftypes.Number, 42),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
			expErrors: 1,
		},
		"other-unknown": {
			req: schemavalidator.NotEqualValidatorRequest{
				ConfigValue:    types.StringValue("bar value"),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.StringAttribute{},
							"bar": schema.StringAttribute{},
							"baz": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.String,
							"bar": tftypes.String,
							"baz": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
						"baz": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := &schemavalidator.NotEqualValidatorResponse{}

			schemavalidator.NotEqualValidator{
				PathExpressions: test.in,
			}.Validate(context.TODO(), test.req, res)

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
