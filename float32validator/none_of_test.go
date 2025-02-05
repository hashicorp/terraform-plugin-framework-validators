// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float32validator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float32validator"
)

func TestNoneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in           types.Float32
		noneOfValues []float32
		expectError  bool
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.Float32Value(123.456),
			noneOfValues: []float32{
				123.456,
				234.567,
				8910.11,
				1213.1415,
			},
			expectError: true,
		},
		"simple-mismatch": {
			in: types.Float32Value(123.456),
			noneOfValues: []float32{
				234.567,
				8910.11,
				1213.1415,
			},
		},
		"skip-validation-on-null": {
			in: types.Float32Null(),
			noneOfValues: []float32{
				234.567,
				8910.11,
				1213.1415,
			},
		},
		"skip-validation-on-unknown": {
			in: types.Float32Unknown(),
			noneOfValues: []float32{
				234.567,
				8910.11,
				1213.1415,
			},
		},
	}

	for name, test := range testCases {
		name, test := name, test

		t.Run(fmt.Sprintf("ValidateFloat32 - %s", name), func(t *testing.T) {
			t.Parallel()
			req := validator.Float32Request{
				ConfigValue: test.in,
			}
			res := validator.Float32Response{}
			float32validator.NoneOf(test.noneOfValues...).ValidateFloat32(context.TODO(), req, &res)

			if !res.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterFloat32 - %s", name), func(t *testing.T) {
			t.Parallel()
			req := function.Float32ParameterValidatorRequest{
				Value: test.in,
			}
			res := function.Float32ParameterValidatorResponse{}
			float32validator.NoneOf(test.noneOfValues...).ValidateParameterFloat32(context.TODO(), req, &res)

			if res.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Error)
			}
		})
	}
}
