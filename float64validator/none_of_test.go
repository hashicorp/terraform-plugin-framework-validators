package float64validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
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

	testCases := map[string]testCase{
		"simple-match": {
			in: types.Float64{Value: 123.456},
			validator: float64validator.NoneOf(
				123.456,
				234.567,
				8910.11,
				1213.1415,
			),
			expErrors: 1,
		},
		"simple-mismatch": {
			in: types.Float64{Value: 123.456},
			validator: float64validator.NoneOf(
				234.567,
				8910.11,
				1213.1415,
			),
			expErrors: 0,
		},
		"skip-validation-on-null": {
			in: types.Float64{Null: true},
			validator: float64validator.NoneOf(
				234.567,
				8910.11,
				1213.1415,
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.Float64{Unknown: true},
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
