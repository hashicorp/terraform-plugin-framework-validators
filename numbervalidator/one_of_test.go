// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package numbervalidator_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
)

func TestOneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in          types.Number
		oneOfValues []*big.Float
		expectError bool
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.NumberValue(big.NewFloat(123.456)),
			oneOfValues: []*big.Float{
				big.NewFloat(123.456),
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
			},
		},
		"simple-mismatch": {
			in: types.NumberValue(big.NewFloat(123.456)),
			oneOfValues: []*big.Float{
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
			},
			expectError: true,
		},
		"skip-validation-on-null": {
			in: types.NumberNull(),
			oneOfValues: []*big.Float{
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
			},
		},
		"skip-validation-on-unknown": {
			in: types.NumberUnknown(),
			oneOfValues: []*big.Float{
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
			},
		},
	}

	for name, test := range testCases {
		name, test := name, test

		t.Run(fmt.Sprintf("ValidateNumber - %s", name), func(t *testing.T) {
			t.Parallel()
			req := validator.NumberRequest{
				ConfigValue: test.in,
			}
			res := validator.NumberResponse{}
			numbervalidator.OneOf(test.oneOfValues...).ValidateNumber(context.TODO(), req, &res)

			if !res.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterNumber - %s", name), func(t *testing.T) {
			t.Parallel()
			req := function.NumberParameterValidatorRequest{
				Value: test.in,
			}
			res := function.NumberParameterValidatorResponse{}
			numbervalidator.OneOf(test.oneOfValues...).ValidateParameterNumber(context.TODO(), req, &res)

			if res.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Error)
			}
		})
	}
}
