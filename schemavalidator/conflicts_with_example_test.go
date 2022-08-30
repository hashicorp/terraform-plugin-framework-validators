package schemavalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleConflictsWith() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					// Validate this attribute must not be configured with other_attr.
					schemavalidator.ConflictsWith(path.Expressions{
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
