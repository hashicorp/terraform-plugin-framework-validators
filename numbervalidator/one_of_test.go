package numbervalidator_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
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

	testCases := map[string]testCase{
		"simple-match": {
			in: types.Number{Value: big.NewFloat(123.456)},
			validator: numbervalidator.OneOf(
				big.NewFloat(123.456),
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
			),
			expErrors: 0,
		},
		"simple-mismatch": {
			in: types.Number{Value: big.NewFloat(123.456)},
			validator: numbervalidator.OneOf(
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
			),
			expErrors: 1,
		},
		"skip-validation-on-null": {
			in: types.Number{Null: true},
			validator: numbervalidator.OneOf(
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.Number{Unknown: true},
			validator: numbervalidator.OneOf(
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
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
