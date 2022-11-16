package stringvalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestOneOfValidator(t *testing.T) {
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
			in: types.StringValue("foo"),
			validator: stringvalidator.OneOf(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"simple-mismatch-case-insensitive": {
			in: types.StringValue("foo"),
			validator: stringvalidator.OneOf(
				"FOO",
				"bar",
				"baz",
			),
			expErrors: 1,
		},
		"simple-mismatch": {
			in: types.StringValue("foz"),
			validator: stringvalidator.OneOf(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 1,
		},
		"list-not-allowed": {
			in: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("10"),
					types.StringValue("20"),
					types.StringValue("30"),
				},
			),
			validator: stringvalidator.OneOf(
				"10",
				"20",
				"30",
				"40",
				"50",
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
			validator: stringvalidator.OneOf(
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
			in: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one.one":    types.StringValue("1.1"),
					"ten.twenty": types.StringValue("10.20"),
					"five.four":  types.StringValue("5.4"),
				},
			),
			validator: stringvalidator.OneOf(
				"1.1",
				"10.20",
				"5.4",
				"geronimo",
				"bob",
			),
			expErrors: 1,
		},
		"object-not-allowed": {
			in: types.ObjectValueMust(
				objAttrTypes,
				map[string]attr.Value{
					"Name":    types.StringValue("Bob Parr"),
					"Age":     types.StringValue("40"),
					"Address": types.StringValue("1200 Park Avenue Emeryville"),
				},
			),
			validator: stringvalidator.OneOf(
				"Bob Parr",
				"40",
				"1200 Park Avenue Emeryville",
				"123",
			),
			expErrors: 1,
		},
		"skip-validation-on-null": {
			in: types.StringNull(),
			validator: stringvalidator.OneOf(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.StringUnknown(),
			validator: stringvalidator.OneOf(
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

func TestOneOfCaseInsensitiveValidator(t *testing.T) {
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
			in: types.StringValue("foo"),
			validator: stringvalidator.OneOfCaseInsensitive(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"simple-match-case-insensitive": {
			in: types.StringValue("foo"),
			validator: stringvalidator.OneOfCaseInsensitive(
				"FOO",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"simple-mismatch": {
			in: types.StringValue("foz"),
			validator: stringvalidator.OneOfCaseInsensitive(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 1,
		},
		"list-not-allowed": {
			in: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("10"),
					types.StringValue("20"),
					types.StringValue("30"),
				},
			),
			validator: stringvalidator.OneOfCaseInsensitive(
				"10",
				"20",
				"30",
				"40",
				"50",
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
			validator: stringvalidator.OneOfCaseInsensitive(
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
			in: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"one.one":    types.StringValue("1.1"),
					"ten.twenty": types.StringValue("10.20"),
					"five.four":  types.StringValue("5.4"),
				},
			),
			validator: stringvalidator.OneOfCaseInsensitive(
				"1.1",
				"10.20",
				"5.4",
				"geronimo",
				"bob",
			),
			expErrors: 1,
		},
		"object-not-allowed": {
			in: types.ObjectValueMust(
				objAttrTypes,
				map[string]attr.Value{
					"Name":    types.StringValue("Bob Parr"),
					"Age":     types.StringValue("40"),
					"Address": types.StringValue("1200 Park Avenue Emeryville"),
				},
			),
			validator: stringvalidator.OneOfCaseInsensitive(
				"Bob Parr",
				"40",
				"1200 Park Avenue Emeryville",
				"123",
			),
			expErrors: 1,
		},
		"skip-validation-on-null": {
			in: types.StringNull(),
			validator: stringvalidator.OneOfCaseInsensitive(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.StringUnknown(),
			validator: stringvalidator.OneOfCaseInsensitive(
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
