package metavalidator_test

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/metavalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleAnyWithAllWarnings() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					// Validate this string value must either be:
					//  - "!!!"
					//  - Any length of alphanumeric characters.
					metavalidator.AnyWithAllWarnings(
						stringvalidator.OneOf("!!!"),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[a-zA-Z0-9]*$`),
							"must contain only alphanumeric characters",
						),
					),
				},
			},
		},
	}
}
