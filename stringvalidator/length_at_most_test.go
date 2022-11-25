package stringvalidator_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestLengthAtMostValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		maxLength   int
		expectError bool
	}
	tests := map[string]testCase{
		"unknown String": {
			val:       types.StringUnknown(),
			maxLength: 1,
		},
		"null String": {
			val:       types.StringNull(),
			maxLength: 1,
		},
		"valid String": {
			val:       types.StringValue("ok"),
			maxLength: 2,
		},
		"too long String": {
			val:         types.StringValue("not ok"),
			maxLength:   5,
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.StringResponse{}
			stringvalidator.LengthAtMost(test.maxLength).ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
