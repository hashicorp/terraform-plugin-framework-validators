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

func TestUTF8LengthAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		maxLength   int
		expectError bool
	}
	tests := map[string]testCase{
		"unknown": {
			val:       types.StringUnknown(),
			maxLength: 1,
		},
		"null": {
			val:       types.StringNull(),
			maxLength: 1,
		},
		"valid single byte characters": {
			val:       types.StringValue("ok"),
			maxLength: 3,
		},
		"valid mixed byte characters": {
			// Rightwards Arrow Over Leftwards Arrow (U+21C4; 3 bytes)
			val:       types.StringValue("test⇄test"),
			maxLength: 9,
		},
		"valid multiple byte characters": {
			// Rightwards Arrow Over Leftwards Arrow (U+21C4; 3 bytes)
			val:       types.StringValue("⇄"),
			maxLength: 1,
		},
		"invalid single byte characters": {
			val:         types.StringValue("ok"),
			maxLength:   1,
			expectError: true,
		},
		"invalid mixed byte characters": {
			// Rightwards Arrow Over Leftwards Arrow (U+21C4; 3 bytes)
			val:         types.StringValue("test⇄test"),
			maxLength:   8,
			expectError: true,
		},
		"invalid multiple byte characters": {
			// Rightwards Arrow Over Leftwards Arrow (U+21C4; 3 bytes)
			val:         types.StringValue("⇄⇄"),
			maxLength:   1,
			expectError: true,
		},
		"invalid validator usage - maxLength < 0": {
			val:         types.StringValue("ok"),
			maxLength:   -1,
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
			stringvalidator.UTF8LengthAtMost(test.maxLength).ValidateString(context.Background(), request, &response)

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
			stringvalidator.UTF8LengthAtMost(test.maxLength).ValidateParameterString(context.Background(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
