package listvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNotNullValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.List
		expectError bool
	}
	tests := map[string]testCase{
		"List null": {
			val: types.ListNull(
				types.StringType,
			),
			expectError: true,
		},
		"List unknown": {
			val: types.ListUnknown(
				types.StringType,
			),
			expectError: false,
		},
		"List empty": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{},
			),
			expectError: false,
		},
		"List with elements": {
			val: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("first"),
				},
			),
			expectError: false,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.ListRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.ListResponse{}
			NotNull().ValidateList(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
