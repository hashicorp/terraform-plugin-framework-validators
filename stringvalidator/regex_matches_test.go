package stringvalidator

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRegexMatchesValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		regexp      *regexp.Regexp
		expectError bool
	}
	tests := map[string]testCase{
		"not a String": {
			val:         types.Bool{Value: true},
			expectError: true,
		},
		"unknown String": {
			val:    types.String{Unknown: true},
			regexp: regexp.MustCompile(`^o[j-l]?$`),
		},
		"null String": {
			val:    types.String{Null: true},
			regexp: regexp.MustCompile(`^o[j-l]?$`),
		},
		"valid String": {
			val:    types.String{Value: "ok"},
			regexp: regexp.MustCompile(`^o[j-l]?$`),
		},
		"invalid String": {
			val:         types.String{Value: "not ok"},
			regexp:      regexp.MustCompile(`^o[j-l]?$`),
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := tfsdk.ValidateAttributeRequest{
				AttributePath:   path.Root("test"),
				AttributeConfig: test.val,
			}
			response := tfsdk.ValidateAttributeResponse{}
			RegexMatches(test.regexp, "").Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
