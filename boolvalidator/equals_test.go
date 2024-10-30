// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boolvalidator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestEqualsValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in          types.Bool
		equalsValue bool
		expectError bool
	}

	testCases := map[string]testCase{
		"simple-match": {
			in:          types.BoolValue(true),
			equalsValue: true,
		},
		"simple-mismatch": {
			in:          types.BoolValue(false),
			equalsValue: true,
			expectError: true,
		},
		"skip-validation-on-null": {
			in:          types.BoolNull(),
			equalsValue: true,
		},
		"skip-validation-on-unknown": {
			in:          types.BoolUnknown(),
			equalsValue: true,
		},
	}

	for name, test := range testCases {
		name, test := name, test

		t.Run(fmt.Sprintf("ValidateBool - %s", name), func(t *testing.T) {
			t.Parallel()
			req := validator.BoolRequest{
				ConfigValue: test.in,
			}
			res := validator.BoolResponse{}
			boolvalidator.Equals(test.equalsValue).ValidateBool(context.TODO(), req, &res)

			if !res.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterBool - %s", name), func(t *testing.T) {
			t.Parallel()
			req := function.BoolParameterValidatorRequest{
				Value: test.in,
			}
			res := function.BoolParameterValidatorResponse{}
			boolvalidator.Equals(test.equalsValue).ValidateParameterBool(context.TODO(), req, &res)

			if res.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Error)
			}
		})
	}
}
