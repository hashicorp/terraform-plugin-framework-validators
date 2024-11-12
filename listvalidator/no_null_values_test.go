// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNoNullValuesValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.List
		expectError bool
	}
	tests := map[string]testCase{
		"List unknown": {
			val: types.ListUnknown(
				types.StringType,
			),
			expectError: false,
		},
		"List null": {
			val: types.ListNull(
				types.StringType,
			),
			expectError: false,
		},
		"No null values": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringValue("second"),
				},
			),
			expectError: false,
		},
		"Unknown value": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringUnknown(),
				},
			),
			expectError: false,
		},
		"Null value": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
					types.StringNull(),
				},
			),
			expectError: true,
		},
	}

	for name, test := range tests {
		t.Run(fmt.Sprintf("ValidateList - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.ListRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.ListResponse{}
			listvalidator.NoNullValues().ValidateList(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterList - %s", name), func(t *testing.T) {
			t.Parallel()
			request := function.ListParameterValidatorRequest{
				ArgumentPosition: 0,
				Value:            test.val,
			}
			response := function.ListParameterValidatorResponse{}
			listvalidator.NoNullValues().ValidateParameterList(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
