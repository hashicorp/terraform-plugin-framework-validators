// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int32validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
)

func TestOneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.Int32
		validator validator.Int32
		expErrors int
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.Int32Value(123),
			validator: int32validator.OneOf(
				123,
				234,
				8910,
				1213,
			),
			expErrors: 0,
		},
		"simple-mismatch": {
			in: types.Int32Value(123),
			validator: int32validator.OneOf(
				234,
				8910,
				1213,
			),
			expErrors: 1,
		},
		"skip-validation-on-null": {
			in: types.Int32Null(),
			validator: int32validator.OneOf(
				234,
				8910,
				1213,
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.Int32Unknown(),
			validator: int32validator.OneOf(
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
			t.Parallel()
			req := validator.Int32Request{
				ConfigValue: test.in,
			}
			res := validator.Int32Response{}
			test.validator.ValidateInt32(context.TODO(), req, &res)

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
