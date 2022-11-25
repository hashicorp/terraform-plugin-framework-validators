package stringvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleAnyWithAllWarnings() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					// Validate this String value must either be:
					//  - "one"
					//  - Length at least 4 characters
					stringvalidator.AnyWithAllWarnings(
						stringvalidator.OneOf("one"),
						stringvalidator.LengthAtLeast(4),
					),
				},
			},
		},
	}
}
