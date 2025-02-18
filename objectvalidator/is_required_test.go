// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package objectvalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
)

func TestIsRequiredValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Object
		expectError bool
	}
	tests := map[string]testCase{
		"Object null": {
			val: types.ObjectNull(
				map[string]attr.Type{
					"field1": types.StringType,
				},
			),
			expectError: true,
		},
		"Object unknown": {
			val: types.ObjectUnknown(
				map[string]attr.Type{
					"field1": types.StringType,
				},
			),
			expectError: false,
		},
		"Object empty": {
			val: types.ObjectValueMust(
				map[string]attr.Type{
					"field1": types.StringType,
				},
				map[string]attr.Value{
					"field1": types.StringNull(),
				},
			),
			expectError: false,
		},
		"Object with elements": {
			val: types.ObjectValueMust(
				map[string]attr.Type{
					"field1": types.StringType,
				},
				map[string]attr.Value{
					"field1": types.StringValue("value1"),
				},
			),
			expectError: false,
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.ObjectRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.ObjectResponse{}
			objectvalidator.IsRequired().ValidateObject(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
