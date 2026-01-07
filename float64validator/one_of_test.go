// Copyright IBM Corp. 2022, 2025
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

func TestOneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in          types.Float64
		oneOfValues []float64
		expectError bool
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.Float64Value(123.456),
			oneOfValues: []float64{
				123.456,
				234.567,
				8910.11,
				1213.1415,
			},
		},
		"simple-mismatch": {
			in: types.Float64Value(123.456),
			oneOfValues: []float64{
				234.567,
				8910.11,
				1213.1415,
			},
			expectError: true,
		},
		"skip-validation-on-null": {
			in: types.Float64Null(),
			oneOfValues: []float64{
				234.567,
				8910.11,
				1213.1415,
			},
		},
		"skip-validation-on-unknown": {
			in: types.Float64Unknown(),
			oneOfValues: []float64{
				234.567,
				8910.11,
				1213.1415,
			},
		},
	}

	for name, test := range testCases {

		t.Run(fmt.Sprintf("ValidateFloat64 - %s", name), func(t *testing.T) {
			t.Parallel()
			req := validator.Float64Request{
				ConfigValue: test.in,
			}
			res := validator.Float64Response{}
			float64validator.OneOf(test.oneOfValues...).ValidateFloat64(context.TODO(), req, &res)

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
			float64validator.OneOf(test.oneOfValues...).ValidateParameterFloat64(context.TODO(), req, &res)

			if res.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Error)
			}
		})
	}
}
