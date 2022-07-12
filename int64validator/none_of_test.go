package int64validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
			in: types.Int64{Value: 123},
			validator: int64validator.NoneOf(
				123,
				234,
				8910,
				1213,
			),
			expErrors: 1,
		},
		"simple-mismatch": {
			in: types.Int64{Value: 123},
			validator: int64validator.NoneOf(
				234,
				8910,
				1213,
			),
			expErrors: 0,
		},
		"skip-validation-on-null": {
			in: types.Int64{Null: true},
			validator: int64validator.NoneOf(
				234,
				8910,
				1213,
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.Int64{Unknown: true},
			validator: int64validator.NoneOf(
				234,
				8910,
				1213,
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
