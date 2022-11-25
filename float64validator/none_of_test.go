package float64validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
)

func TestNoneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.Float64
		validator validator.Float64
		expErrors int
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.Float64Value(123.456),
			validator: float64validator.NoneOf(
				123.456,
				234.567,
				8910.11,
				1213.1415,
			),
			expErrors: 1,
		},
		"simple-mismatch": {
			in: types.Float64Value(123.456),
			validator: float64validator.NoneOf(
				234.567,
				8910.11,
				1213.1415,
			),
			expErrors: 0,
		},
		"skip-validation-on-null": {
			in: types.Float64Null(),
			validator: float64validator.NoneOf(
				234.567,
				8910.11,
				1213.1415,
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.Float64Unknown(),
			validator: float64validator.NoneOf(
				234.567,
				8910.11,
				1213.1415,
			),
			expErrors: 0,
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			req := validator.Float64Request{
				ConfigValue: test.in,
			}
			res := validator.Float64Response{}
			test.validator.ValidateFloat64(context.TODO(), req, &res)

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
