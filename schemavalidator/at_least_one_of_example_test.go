package schemavalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleAtLeastOneOf() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					// Validate at least this attribute or other_attr should be configured.
					schemavalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("other_attr"),
					}...),
				},
			},
			"other_attr": {
				Required: true,
				Type:     types.StringType,
			},
		},
	}
}
