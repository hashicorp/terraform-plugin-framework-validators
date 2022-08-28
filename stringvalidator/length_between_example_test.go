package stringvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleLengthBetween() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					// Validate string value length must be at least 3 and at most 256 characters.
					stringvalidator.LengthBetween(3, 256),
				},
			},
		},
	}
}
