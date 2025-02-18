// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestLengthAtLeastValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		minLength   int
		expectError bool
	}
	tests := map[string]testCase{
		"unknown": {
			val:       types.StringUnknown(),
			minLength: 1,
		},
		"null": {
			val:       types.StringNull(),
			minLength: 1,
		},
		"valid": {
			val:       types.StringValue("ok"),
			minLength: 1,
		},
		"too short": {
			val:         types.StringValue(""),
			minLength:   1,
			expectError: true,
		},
		"multiple byte characters": {
			// Rightwards Arrow Over Leftwards Arrow (U+21C4; 3 bytes)
			val:       types.StringValue("⇄"),
			minLength: 2,
		},
		"invalid validator usage - minLength < 0": {
			val:         types.StringValue("ok"),
			minLength:   -1,
			expectError: true,
		},
	}

	for name, test := range tests {

		t.Run(fmt.Sprintf("ValidateString - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.StringResponse{}
			stringvalidator.LengthAtLeast(test.minLength).ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterString - %s", name), func(t *testing.T) {
			t.Parallel()
			request := function.StringParameterValidatorRequest{
				ArgumentPosition: 0,
				Value:            test.val,
			}
			response := function.StringParameterValidatorResponse{}
			stringvalidator.LengthAtLeast(test.minLength).ValidateParameterString(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
