package primitivevalidator_test

import (
	"context"
	"math"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/primitivevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestOneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        attr.Value
		validator tfsdk.AttributeValidator
		expErrors int
	}

	objPersonAttrTypes := map[string]attr.Type{
		"Name": types.StringType,
		"Age":  types.Int64Type,
	}
	objAttrTypes := map[string]attr.Type{
		"Person": types.ObjectType{
			AttrTypes: objPersonAttrTypes,
		},
		"Address": types.StringType,
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.String{Value: "foo"},
			validator: primitivevalidator.OneOf(
				types.String{Value: "foo"},
				types.String{Value: "bar"},
				types.String{Value: "baz"},
			),
			expErrors: 0,
		},
		"simple-mismatch": {
			in: types.String{Value: "foz"},
			validator: primitivevalidator.OneOf(
				types.String{Value: "foo"},
				types.String{Value: "bar"},
				types.String{Value: "baz"},
			),
			expErrors: 1,
		},
		"mixed": {
			in: types.Float64{Value: 1.234},
			validator: primitivevalidator.OneOf(
				types.String{Value: "foo"},
				types.Int64{Value: 567},
				types.Float64{Value: 1.234},
			),
			expErrors: 0,
		},
		"list-not-allowed": {
			in: types.List{
				ElemType: types.Int64Type,
				Elems: []attr.Value{
					types.Int64{Value: 10},
					types.Int64{Value: 20},
					types.Int64{Value: 30},
				},
			},
			validator: primitivevalidator.OneOf(
				types.Int64{Value: 10},
				types.Int64{Value: 20},
				types.Int64{Value: 30},
				types.Int64{Value: 40},
				types.Int64{Value: 50},
			),
			expErrors: 1,
		},
		"set-not-allowed": {
			in: types.Set{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "foo"},
					types.String{Value: "bar"},
					types.String{Value: "baz"},
				},
			},
			validator: primitivevalidator.OneOf(
				types.String{Value: "bob"},
				types.String{Value: "alice"},
				types.String{Value: "john"},
				types.String{Value: "foo"},
				types.String{Value: "bar"},
				types.String{Value: "baz"},
			),
			expErrors: 1,
		},
		"map-not-allowed": {
			in: types.Map{
				ElemType: types.NumberType,
				Elems: map[string]attr.Value{
					"one.one":    types.Number{Value: big.NewFloat(1.1)},
					"ten.twenty": types.Number{Value: big.NewFloat(10.20)},
					"five.four":  types.Number{Value: big.NewFloat(5.4)},
				},
			},
			validator: primitivevalidator.OneOf(
				types.Number{Value: big.NewFloat(1.1)},
				types.Number{Value: big.NewFloat(math.MaxFloat64)},
				types.Number{Value: big.NewFloat(math.SmallestNonzeroFloat64)},
				types.Number{Value: big.NewFloat(10.20)},
				types.Number{Value: big.NewFloat(5.4)},
			),
			expErrors: 1,
		},
		"object-not-allowed": {
			in: types.Object{
				AttrTypes: objAttrTypes,
				Attrs: map[string]attr.Value{
					"Person": types.Object{
						AttrTypes: objPersonAttrTypes,
						Attrs: map[string]attr.Value{
							"Name": types.String{Value: "Bob Parr"},
							"Age":  types.Int64{Value: 40},
						},
					},
					"Address": types.String{Value: "1200 Park Avenue Emeryville"},
				},
			},
			validator: primitivevalidator.OneOf(
				types.Object{
					AttrTypes: map[string]attr.Type{},
					Attrs:     map[string]attr.Value{},
				},
				types.Object{
					AttrTypes: objPersonAttrTypes,
					Attrs: map[string]attr.Value{
						"Name": types.String{Value: "Bob Parr"},
						"Age":  types.Int64{Value: 40},
					},
				},
				types.String{Value: "1200 Park Avenue Emeryville"},
				types.Int64{Value: 123},
				types.String{Value: "Bob Parr"},
			),
			expErrors: 1,
		},
		"skip-validation-on-null": {
			in: types.String{Null: true},
			validator: primitivevalidator.OneOf(
				types.String{Value: "foo"},
				types.String{Value: "bar"},
				types.String{Value: "baz"},
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.String{Unknown: true},
			validator: primitivevalidator.OneOf(
				types.String{Value: "foo"},
				types.String{Value: "bar"},
				types.String{Value: "baz"},
			),
			expErrors: 0,
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			req := tfsdk.ValidateAttributeRequest{
				AttributeConfig: test.in,
			}
			res := tfsdk.ValidateAttributeResponse{}
			test.validator.Validate(context.TODO(), req, &res)

			if test.expErrors > 0 && !res.Diagnostics.HasError() {
				t.Fatalf("expected %d error(s), got none", test.expErrors)
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
