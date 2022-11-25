package stringvalidator_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestRegexMatchesValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		regexp      *regexp.Regexp
		expectError bool
	}
	tests := map[string]testCase{
		"unknown String": {
			val:    types.StringUnknown(),
			regexp: regexp.MustCompile(`^o[j-l]?$`),
		},
		"null String": {
			val:    types.StringNull(),
			regexp: regexp.MustCompile(`^o[j-l]?$`),
		},
		"valid String": {
			val:    types.StringValue("ok"),
			regexp: regexp.MustCompile(`^o[j-l]?$`),
		},
		"invalid String": {
			val:         types.StringValue("not ok"),
			regexp:      regexp.MustCompile(`^o[j-l]?$`),
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
			stringvalidator.RegexMatches(test.regexp, "").ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
