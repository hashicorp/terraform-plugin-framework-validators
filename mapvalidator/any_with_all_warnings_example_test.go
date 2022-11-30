package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleAnyWithAllWarnings() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.MapAttribute{
				Required: true,
				Validators: []validator.Map{
					// Validate this Map value must either be:
					//  - Between 1 and 2 elements
					//  - At least 4 elements
					mapvalidator.AnyWithAllWarnings(
						mapvalidator.SizeBetween(1, 2),
						mapvalidator.SizeAtLeast(4),
					),
				},
			},
		},
	}
}
