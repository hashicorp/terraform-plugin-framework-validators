package primitivevalidator_test

import (
	"context"
	"math"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/primitivevalidator"
)

func TestNoneOfValidator(t *testing.T) {
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
			in: types.StringValue("foo"),
			validator: primitivevalidator.NoneOf(
				types.StringValue("foo"),
				types.StringValue("bar"),
				types.StringValue("baz"),
			),
			expErrors: 1,
		},
		"simple-mismatch": {
			in: types.StringValue("foz"),
			validator: primitivevalidator.NoneOf(
				types.StringValue("foo"),
				types.StringValue("bar"),
				types.StringValue("baz"),
			),
			expErrors: 0,
		},
		"mixed": {
			in: types.Float64Value(1.234),
			validator: primitivevalidator.NoneOf(
				types.StringValue("foo"),
				types.Int64Value(567),
				types.Float64Value(1.234),
			),
			expErrors: 1,
		},
		"list-not-allowed": {
			in: types.ListValueMust(
				types.Int64Type,
				[]attr.Value{
					types.Int64Value(10),
					types.Int64Value(20),
					types.Int64Value(30),
				},
			),
			validator: primitivevalidator.NoneOf(
				types.Int64Value(10),
				types.Int64Value(20),
				types.Int64Value(30),
				types.Int64Value(40),
				types.Int64Value(50),
			),
			expErrors: 1,
		},
		"set-not-allowed": {
			in: types.SetValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("foo"),
					types.StringValue("bar"),
					types.StringValue("baz"),
				},
			),
			validator: primitivevalidator.NoneOf(
				types.StringValue("bob"),
				types.StringValue("alice"),
				types.StringValue("john"),
				types.StringValue("foo"),
				types.StringValue("bar"),
				types.StringValue("baz"),
			),
			expErrors: 1,
		},
		"map-not-allowed": {
			in: types.MapValueMust(
				types.NumberType,
				map[string]attr.Value{
					"one.one":    types.NumberValue(big.NewFloat(1.1)),
					"ten.twenty": types.NumberValue(big.NewFloat(10.20)),
					"five.four":  types.NumberValue(big.NewFloat(5.4)),
				},
			),
			validator: primitivevalidator.NoneOf(
				types.NumberValue(big.NewFloat(1.1)),
				types.NumberValue(big.NewFloat(math.MaxFloat64)),
				types.NumberValue(big.NewFloat(math.SmallestNonzeroFloat64)),
				types.NumberValue(big.NewFloat(10.20)),
				types.NumberValue(big.NewFloat(5.4)),
			),
			expErrors: 1,
		},
		"object-not-allowed": {
			in: types.ObjectValueMust(
				objAttrTypes,
				map[string]attr.Value{
					"Person": types.ObjectValueMust(
						objPersonAttrTypes,
						map[string]attr.Value{
							"Name": types.StringValue("Bob Parr"),
							"Age":  types.Int64Value(40),
						},
					),
					"Address": types.StringValue("1200 Park Avenue Emeryville"),
				},
			),
			validator: primitivevalidator.NoneOf(
				types.ObjectValueMust(
					map[string]attr.Type{},
					map[string]attr.Value{},
				),
				types.ObjectValueMust(
					objPersonAttrTypes,
					map[string]attr.Value{
						"Name": types.StringValue("Bob Parr"),
						"Age":  types.Int64Value(40),
					},
				),
				types.StringValue("1200 Park Avenue Emeryville"),
				types.Int64Value(123),
				types.StringValue("Bob Parr"),
			),
			expErrors: 1,
		},
		"skip-validation-on-null": {
			in: types.StringNull(),
			validator: primitivevalidator.NoneOf(
				types.StringValue("foo"),
				types.StringValue("bar"),
				types.StringValue("baz"),
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.StringUnknown(),
			validator: primitivevalidator.NoneOf(
				types.StringValue("foo"),
				types.StringValue("bar"),
				types.StringValue("baz"),
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

			if test.expErrors > 0 && test.expErrors != res.Diagnostics.ErrorsCount() {
				t.Fatalf("expected %d error(s), got %d: %v", test.expErrors, res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}

			if test.expErrors == 0 && res.Diagnostics.HasError() {
				t.Fatalf("expected no error(s), got %d: %v", res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}
		})
	}
}
