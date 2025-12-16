// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package mapvalidator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNoNullValuesValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Map
		expectError bool
	}
	tests := map[string]testCase{
		"Map unknown": {
			val: types.MapUnknown(
				types.StringType,
			),
			expectError: false,
		},
		"Map null": {
			val: types.MapNull(
				types.StringType,
			),
			expectError: false,
		},
		"No null values": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("first"),
					"key2": types.StringValue("second"),
				},
			),
			expectError: false,
		},
		"Unknown value": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("first"),
					"key2": types.StringUnknown(),
				},
			),
			expectError: false,
		},
		"Null value": {
			val: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("first"),
					"key2": types.StringNull(),
				},
			),
			expectError: true,
		},
	}

	for name, test := range tests {
		t.Run(fmt.Sprintf("ValidateMap - %s", name), func(t *testing.T) {
			t.Parallel()
			request := validator.MapRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.MapResponse{}
			mapvalidator.NoNullValues().ValidateMap(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterMap - %s", name), func(t *testing.T) {
			t.Parallel()
			request := function.MapParameterValidatorRequest{
				ArgumentPosition: 0,
				Value:            test.val,
			}
			response := function.MapParameterValidatorResponse{}
			mapvalidator.NoNullValues().ValidateParameterMap(context.TODO(), request, &response)

			if response.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Error)
			}
		})
	}
}
