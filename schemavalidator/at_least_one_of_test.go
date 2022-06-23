package schemavalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAtLeastOneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		req       tfsdk.ValidateAttributeRequest
		in        []*tftypes.AttributePath
		expErrors int
	}

	testCases := map[string]testCase{
		"base": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.String{Value: "bar value"},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("bar"),
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
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
					}),
				},
			},
			in: []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("foo"),
			},
		},
		"self-is-null": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.String{Null: true},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("bar"),
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
			in: []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("foo"),
			},
		},
		"error_none-set": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.String{Value: "bar value"},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("bar"),
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
						"foo": tftypes.NewValue(tftypes.Number, nil),
						"bar": tftypes.NewValue(tftypes.String, nil),
						"baz": tftypes.NewValue(tftypes.Number, nil),
					}),
				},
			},
			in: []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("foo"),
				tftypes.NewAttributePath().WithAttributeName("baz"),
			},
			expErrors: 1,
		},
		"multiple-set": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.String{Value: "bar value"},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("bar"),
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
								Type: types.Float64Type,
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
						"baz": tftypes.NewValue(tftypes.Number, 4.2),
					}),
				},
			},
			in: []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("foo"),
				tftypes.NewAttributePath().WithAttributeName("baz"),
			},
		},
		"allow-duplicate-input": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.String{Value: "bar value"},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("bar"),
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
						"foo": tftypes.NewValue(tftypes.Number, nil),
						"bar": tftypes.NewValue(tftypes.String, "bar value"),
						"baz": tftypes.NewValue(tftypes.Number, nil),
					}),
				},
			},
			in: []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("foo"),
				tftypes.NewAttributePath().WithAttributeName("bar"),
				tftypes.NewAttributePath().WithAttributeName("baz"),
			},
		},
		"unknowns": {
			req: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.String{Value: "bar value"},
				AttributePath:   tftypes.NewAttributePath().WithAttributeName("bar"),
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
			in: []*tftypes.AttributePath{
				tftypes.NewAttributePath().WithAttributeName("foo"),
				tftypes.NewAttributePath().WithAttributeName("baz"),
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			res := tfsdk.ValidateAttributeResponse{}

			schemavalidator.AtLeastOneOf(test.in...).Validate(context.TODO(), test.req, &res)

			if test.expErrors > 0 && !res.Diagnostics.HasError() {
				t.Fatal("expected error(s), got none")
			}

			if test.expErrors > 0 && test.expErrors != validatordiag.ErrorsCount(res.Diagnostics) {
				t.Fatalf("expected %d error(s), got %d: %v", test.expErrors, validatordiag.ErrorsCount(res.Diagnostics), res.Diagnostics)
			}

			if test.expErrors == 0 && res.Diagnostics.HasError() {
				t.Fatalf("expected no error(s), got %d: %v", validatordiag.ErrorsCount(res.Diagnostics), res.Diagnostics)
			}
		})
	}
}
