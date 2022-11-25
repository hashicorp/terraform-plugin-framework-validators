package stringvalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNoneOfCaseInsensitiveValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.String
		validator validator.String
		expErrors int
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.StringValue("foo"),
			validator: stringvalidator.NoneOfCaseInsensitive(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 1,
		},
		"simple-match-case-insensitive": {
			in: types.StringValue("foo"),
			validator: stringvalidator.NoneOfCaseInsensitive(
				"FOO",
				"bar",
				"baz",
			),
			expErrors: 1,
		},
		"simple-mismatch": {
			in: types.StringValue("foz"),
			validator: stringvalidator.NoneOfCaseInsensitive(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"skip-validation-on-null": {
			in: types.StringNull(),
			validator: stringvalidator.NoneOfCaseInsensitive(
				"foo",
				"bar",
				"baz",
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.StringUnknown(),
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
			req := validator.StringRequest{
				ConfigValue: test.in,
			}
			res := validator.StringResponse{}
			test.validator.ValidateString(context.TODO(), req, &res)

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
