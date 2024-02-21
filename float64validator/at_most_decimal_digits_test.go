// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float64validator_test

import (
	"context"
	"math"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
)

func TestAtMostDecimalDigitsValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val                 types.Float64
		atMostDecimalDigits int
		expectError         bool
	}
	tests := map[string]testCase{
		"unknown Float64": {
			val:                 types.Float64Unknown(),
			atMostDecimalDigits: 1,
		},
		"null Float64": {
			val:                 types.Float64Null(),
			atMostDecimalDigits: 2,
		},
		"valid integer as Float64": {
			val:                 types.Float64Value(1),
			atMostDecimalDigits: 2,
		},
		"valid float as Float64": {
			val:                 types.Float64Value(1.1),
			atMostDecimalDigits: 2,
		},
		"valid float as Float64 max": {
			val:                 types.Float64Value(2.0),
			atMostDecimalDigits: 2,
		},
		"too large float as Float64": {
			val:                 types.Float64Value(math.MaxFloat64),
			atMostDecimalDigits: 2,
		},
		"zero decimal digits": {
			val:                 types.Float64Value(3.0000),
			atMostDecimalDigits: 2,
		},
		"more than allowed decimal digits": {
			val:                 types.Float64Value(3.00099),
			atMostDecimalDigits: 3,
			expectError:         true,
		},
		"exactly same as allowed decimal digits": {
			val:                 types.Float64Value(54545.009),
			atMostDecimalDigits: 3,
		},
		"less than allowed decimal digits": {
			val:                 types.Float64Value(0.09),
			atMostDecimalDigits: 3,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.Float64Request{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.Float64Response{}
			float64validator.AtMostDecimalDigits(test.atMostDecimalDigits).ValidateFloat64(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
