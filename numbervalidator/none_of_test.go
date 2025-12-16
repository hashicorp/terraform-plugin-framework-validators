// Copyright IBM Corp. 2022, 2025
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

func TestNoneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in           types.Number
		noneOfValues []*big.Float
		expectError  bool
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.NumberValue(big.NewFloat(123.456)),
			noneOfValues: []*big.Float{
				big.NewFloat(123.456),
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
			},
			expectError: true,
		},
		"simple-mismatch": {
			in: types.NumberValue(big.NewFloat(123.456)),
			noneOfValues: []*big.Float{
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
			},
		},
		"skip-validation-on-null": {
			in: types.NumberNull(),
			noneOfValues: []*big.Float{
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
			},
		},
		"skip-validation-on-unknown": {
			in: types.NumberUnknown(),
			noneOfValues: []*big.Float{
				big.NewFloat(234.567),
				big.NewFloat(8910.11),
				big.NewFloat(1213.1415),
			},
		},
	}

	for name, test := range testCases {

		t.Run(fmt.Sprintf("ValidateNumber - %s", name), func(t *testing.T) {
			t.Parallel()
			req := validator.NumberRequest{
				ConfigValue: test.in,
			}
			res := validator.NumberResponse{}
			numbervalidator.NoneOf(test.noneOfValues...).ValidateNumber(context.TODO(), req, &res)

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
			numbervalidator.NoneOf(test.noneOfValues...).ValidateParameterNumber(context.TODO(), req, &res)

			if res.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Error)
			}
		})
	}
}
