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

func TestIsJsonValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		expectError bool
	}
	tests := map[string]testCase{
		"unknown": {
			val: types.StringUnknown(),
		},
		"null": {
			val: types.StringNull(),
		},
		"empty": {
			val: types.StringValue(``),
		},
		"empty brackets": {
			val: types.StringValue(`{}`),
		},
		"valid json": {
			val: types.StringValue(`{"abc":["1","2"]}`),
		},
		"Invalid 1": {
			val:         types.StringValue(`{0:"1"}`),
			expectError: true,
		},
		"Invalid 2": {
			val:         types.StringValue(`{'abc':1}`),
			expectError: true,
		},
		"Invalid 3": {
			val:         types.StringValue(`{"def":}`),
			expectError: true,
		},
		"Invalid 4": {
			val:         types.StringValue(`{"xyz":[}}`),
			expectError: true,
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
			stringvalidator.IsJson().ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
