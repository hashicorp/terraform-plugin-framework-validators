// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float32validator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float32validator"
)

func TestOneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.Float32
		validator validator.Float32
		expErrors int
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.Float32Value(123.456),
			validator: float32validator.OneOf(
				123.456,
				234.567,
				8910.11,
				1213.1415,
			),
			expErrors: 0,
		},
		"simple-mismatch": {
			in: types.Float32Value(123.456),
			validator: float32validator.OneOf(
				234.567,
				8910.11,
				1213.1415,
			),
			expErrors: 1,
		},
		"skip-validation-on-null": {
			in: types.Float32Null(),
			validator: float32validator.OneOf(
				234.567,
				8910.11,
				1213.1415,
			),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in: types.Float32Unknown(),
			validator: float32validator.OneOf(
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
			t.Parallel()
			req := validator.Float32Request{
				ConfigValue: test.in,
			}
			res := validator.Float32Response{}
			test.validator.ValidateFloat32(context.TODO(), req, &res)

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
