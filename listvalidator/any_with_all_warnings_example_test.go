package listvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleAnyWithAllWarnings() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.ListAttribute{
				Required: true,
				Validators: []validator.List{
					// Validate this List value must either be:
					//  - Between 1 and 2 elements
					//  - At least 4 elements
					listvalidator.AnyWithAllWarnings(
						listvalidator.SizeBetween(1, 2),
						listvalidator.SizeAtLeast(4),
					),
				},
			},
		},
	}
}
