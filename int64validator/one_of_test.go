// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int64validator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
)

func TestOneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in          types.Int64
		oneOfValues []int64
		expectError bool
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.Int64Value(123),
			oneOfValues: []int64{
				123,
				234,
				8910,
				1213,
			},
		},
		"simple-mismatch": {
			in: types.Int64Value(123),
			oneOfValues: []int64{
				234,
				8910,
				1213,
			},
			expectError: true,
		},
		"skip-validation-on-null": {
			in: types.Int64Null(),
			oneOfValues: []int64{
				234,
				8910,
				1213,
			},
		},
		"skip-validation-on-unknown": {
			in: types.Int64Unknown(),
			oneOfValues: []int64{
				234,
				8910,
				1213,
			},
		},
	}

	for name, test := range testCases {

		t.Run(fmt.Sprintf("ValidateInt64 - %s", name), func(t *testing.T) {
			t.Parallel()
			req := validator.Int64Request{
				ConfigValue: test.in,
			}
			res := validator.Int64Response{}
			int64validator.OneOf(test.oneOfValues...).ValidateInt64(context.TODO(), req, &res)

			if !res.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterInt64 - %s", name), func(t *testing.T) {
			t.Parallel()
			req := function.Int64ParameterValidatorRequest{
				Value: test.in,
			}
			res := function.Int64ParameterValidatorResponse{}
			int64validator.OneOf(test.oneOfValues...).ValidateParameterInt64(context.TODO(), req, &res)

			if res.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Error)
			}
		})
	}
}
