package stringvalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNoneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        attr.Value
		validator tfsdk.AttributeValidator
		expErrors int
	}

	objAttrTypes := map[string]attr.Type{
		"Name":    types.StringType,
		"Age":     types.StringType,
		"Address": types.StringType,
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.String{Value: "foo"},
			validator: stringvalidator.NoneOf(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 1,
		},
		"simple-mismatch-case-insensitive": {
			in: types.String{Value: "foo"},
			validator: stringvalidator.NoneOf(
				"FOO",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"simple-mismatch": {
			in: types.String{Value: "foz"},
			validator: stringvalidator.NoneOf(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"list-not-allowed": {
			in: types.List{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "10"},
					types.String{Value: "20"},
					types.String{Value: "30"},
				},
			},
			validator: stringvalidator.NoneOf(
				"10",
				"20",
				"30",
				"40",
				"50",
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
			validator: stringvalidator.NoneOf(
				"bob",
				"alice",
				"john",
				"foo",
				"bar",
				"baz",
			),
			expErrors: 1,
		},
		"map-not-allowed": {
			in: types.Map{
				ElemType: types.StringType,
				Elems: map[string]attr.Value{
					"one.one":    types.String{Value: "1.1"},
					"ten.twenty": types.String{Value: "10.20"},
					"five.four":  types.String{Value: "5.4"},
				},
			},
			validator: stringvalidator.NoneOf(
				"1.1",
				"10.20",
				"5.4",
				"geronimo",
				"bob",
			),
			expErrors: 1,
		},
		"object-not-allowed": {
			in: types.Object{
				AttrTypes: objAttrTypes,
				Attrs: map[string]attr.Value{
					"Name":    types.String{Value: "Bob Parr"},
					"Age":     types.String{Value: "40"},
					"Address": types.String{Value: "1200 Park Avenue Emeryville"},
				},
			},
			validator: stringvalidator.NoneOf(
				"Bob Parr",
				"40",
				"1200 Park Avenue Emeryville",
				"123",
			),
			expErrors: 1,
		},
		"skip-validation-on-null": {
			in: types.String{Null: true},
			validator: stringvalidator.NoneOf(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.String{Unknown: true},
			validator: stringvalidator.NoneOf(
				"foo",
				"bar",
				"baz",
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

func TestNoneOfCaseInsensitiveValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        attr.Value
		validator tfsdk.AttributeValidator
		expErrors int
	}

	objAttrTypes := map[string]attr.Type{
		"Name":    types.StringType,
		"Age":     types.StringType,
		"Address": types.StringType,
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.String{Value: "foo"},
			validator: stringvalidator.NoneOfCaseInsensitive(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 1,
		},
		"simple-match-case-insensitive": {
			in: types.String{Value: "foo"},
			validator: stringvalidator.NoneOfCaseInsensitive(
				"FOO",
				"bar",
				"baz",
			),
			expErrors: 1,
		},
		"simple-mismatch": {
			in: types.String{Value: "foz"},
			validator: stringvalidator.NoneOfCaseInsensitive(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"list-not-allowed": {
			in: types.List{
				ElemType: types.StringType,
				Elems: []attr.Value{
					types.String{Value: "10"},
					types.String{Value: "20"},
					types.String{Value: "30"},
				},
			},
			validator: stringvalidator.NoneOfCaseInsensitive(
				"10",
				"20",
				"30",
				"40",
				"50",
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
			validator: stringvalidator.NoneOfCaseInsensitive(
				"bob",
				"alice",
				"john",
				"foo",
				"bar",
				"baz",
			),
			expErrors: 1,
		},
		"map-not-allowed": {
			in: types.Map{
				ElemType: types.StringType,
				Elems: map[string]attr.Value{
					"one.one":    types.String{Value: "1.1"},
					"ten.twenty": types.String{Value: "10.20"},
					"five.four":  types.String{Value: "5.4"},
				},
			},
			validator: stringvalidator.NoneOfCaseInsensitive(
				"1.1",
				"10.20",
				"5.4",
				"geronimo",
				"bob",
			),
			expErrors: 1,
		},
		"object-not-allowed": {
			in: types.Object{
				AttrTypes: objAttrTypes,
				Attrs: map[string]attr.Value{
					"Name":    types.String{Value: "Bob Parr"},
					"Age":     types.String{Value: "40"},
					"Address": types.String{Value: "1200 Park Avenue Emeryville"},
				},
			},
			validator: stringvalidator.NoneOfCaseInsensitive(
				"Bob Parr",
				"40",
				"1200 Park Avenue Emeryville",
				"123",
			),
			expErrors: 1,
		},
		"skip-validation-on-null": {
			in: types.String{Null: true},
			validator: stringvalidator.NoneOfCaseInsensitive(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.String{Unknown: true},
			validator: stringvalidator.NoneOfCaseInsensitive(
				"foo",
				"bar",
				"baz",
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
