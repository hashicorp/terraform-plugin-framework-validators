package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleAll() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Map{
					// Validate this Map value must either be:
					//  - More than 5 elements
					//  - At least 2 elements, but not more than 3 elements
					mapvalidator.Any(
						mapvalidator.SizeAtLeast(5),
						mapvalidator.All(
							mapvalidator.SizeAtLeast(2),
							mapvalidator.SizeAtMost(3),
						),
					),
				},
			},
		},
	}
}
