// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float64validator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
)

func TestNoneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in           types.Float64
		noneOfValues []float64
		expectError  bool
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.Float64Value(123.456),
			noneOfValues: []float64{
				123.456,
				234.567,
				8910.11,
				1213.1415,
			},
			expectError: true,
		},
		"simple-mismatch": {
			in: types.Float64Value(123.456),
			noneOfValues: []float64{
				234.567,
				8910.11,
				1213.1415,
			},
		},
		"skip-validation-on-null": {
			in: types.Float64Null(),
			noneOfValues: []float64{
				234.567,
				8910.11,
				1213.1415,
			},
		},
		"skip-validation-on-unknown": {
			in: types.Float64Unknown(),
			noneOfValues: []float64{
				234.567,
				8910.11,
				1213.1415,
			},
		},
	}

	for name, test := range testCases {
		name, test := name, test

		t.Run(fmt.Sprintf("ValidateFloat64 - %s", name), func(t *testing.T) {
			t.Parallel()
			req := validator.Float64Request{
				ConfigValue: test.in,
			}
			res := validator.Float64Response{}
			float64validator.NoneOf(test.noneOfValues...).ValidateFloat64(context.TODO(), req, &res)

			if !res.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterFloat64 - %s", name), func(t *testing.T) {
			t.Parallel()
			req := function.Float64ParameterValidatorRequest{
				Value: test.in,
			}
			res := function.Float64ParameterValidatorResponse{}
			float64validator.NoneOf(test.noneOfValues...).ValidateParameterFloat64(context.TODO(), req, &res)

			if res.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Error)
			}
		})
	}
}
