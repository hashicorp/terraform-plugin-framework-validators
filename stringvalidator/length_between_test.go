// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestLengthBetweenValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		minLength   int
		maxLength   int
		expectError bool
	}
	tests := map[string]testCase{
		"unknown": {
			val:       types.StringUnknown(),
			minLength: 1,
			maxLength: 3,
		},
		"null": {
			val:       types.StringNull(),
			minLength: 1,
			maxLength: 3,
		},
		"valid": {
			val:       types.StringValue("ok"),
			minLength: 1,
			maxLength: 3,
		},
		"too long": {
			val:         types.StringValue("not ok"),
			minLength:   1,
			maxLength:   3,
			expectError: true,
		},
		"too short": {
			val:         types.StringValue(""),
			minLength:   1,
			maxLength:   3,
			expectError: true,
		},
		"multiple byte characters": {
			// Rightwards Arrow Over Leftwards Arrow (U+21C4; 3 bytes)
			val:       types.StringValue("⇄"),
			minLength: 2,
			maxLength: 4,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.StringResponse{}
			stringvalidator.LengthBetween(test.minLength, test.maxLength).ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
