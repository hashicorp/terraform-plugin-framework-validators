// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemavalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator"
)

func TestAtLeastOneOfValidatorValidate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		req       schemavalidator.AtLeastOneOfValidatorRequest
		in        path.Expressions
		expErrors int
	}

	testCases := map[string]testCase{
		"base": {
			req: schemavalidator.AtLeastOneOfValidatorRequest{
				ConfigValue:    types.StringValue("bar value"),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.Int64Attribute{},
							"bar": schema.StringAttribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, 42),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
			},
		},
		"self-is-null": {
			req: schemavalidator.AtLeastOneOfValidatorRequest{
				ConfigValue:    types.StringNull(),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.Int64Attribute{},
							"bar": schema.StringAttribute{},
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
		"error_none-set": {
			req: schemavalidator.AtLeastOneOfValidatorRequest{
				ConfigValue:    types.StringNull(),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.Int64Attribute{},
							"bar": schema.StringAttribute{},
							"baz": schema.Int64Attribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.String,
							"baz": tftypes.Number,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, nil),
						"bar": tftypes.NewValue(tftypes.String, nil),
						"baz": tftypes.NewValue(tftypes.Number, nil),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
				path.MatchRoot("baz"),
			},
			expErrors: 1,
		},
		"multiple-set": {
			req: schemavalidator.AtLeastOneOfValidatorRequest{
				ConfigValue:    types.StringValue("bar value"),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.Int64Attribute{},
							"bar": schema.StringAttribute{},
							"baz": schema.Float64Attribute{},
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
						"baz": tftypes.NewValue(tftypes.Number, 4.2),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
				path.MatchRoot("baz"),
			},
		},
		"allow-duplicate-input": {
			req: schemavalidator.AtLeastOneOfValidatorRequest{
				ConfigValue:    types.StringValue("bar value"),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.Int64Attribute{},
							"bar": schema.StringAttribute{},
							"baz": schema.Int64Attribute{},
						},
					},
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"foo": tftypes.Number,
							"bar": tftypes.String,
							"baz": tftypes.Number,
						},
					}, map[string]tftypes.Value{
						"foo": tftypes.NewValue(tftypes.Number, nil),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
						"baz": tftypes.NewValue(tftypes.Number, nil),
					}),
				},
			},
			in: path.Expressions{
				path.MatchRoot("foo"),
				path.MatchRoot("bar"),
				path.MatchRoot("baz"),
			},
		},
		"unknowns": {
			req: schemavalidator.AtLeastOneOfValidatorRequest{
				ConfigValue:    types.StringValue("bar value"),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.Int64Attribute{},
							"bar": schema.StringAttribute{},
							"baz": schema.Int64Attribute{},
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
		},
		"matches-no-attribute-in-schema": {
			req: schemavalidator.AtLeastOneOfValidatorRequest{
				ConfigValue:    types.StringNull(),
				Path:           path.Root("bar"),
				PathExpression: path.MatchRoot("bar"),
				Config: tfsdk.Config{
					Schema: schema.Schema{
						Attributes: map[string]schema.Attribute{
							"foo": schema.Int64Attribute{},
							"bar": schema.StringAttribute{},
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
				path.MatchRoot("fooz"),
			},
			expErrors: 2,
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := &schemavalidator.AtLeastOneOfValidatorResponse{}

			schemavalidator.AtLeastOneOfValidator{
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
