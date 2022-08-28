package stringvalidator_test

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleRegexMatches() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					// Validate string value satisfies the regular expression for alphanumeric characters
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9]*$`),
						"must only contain only alphanumeric characters",
					),
				},
			},
		},
	}
}
