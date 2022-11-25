package numbervalidator_test

import (
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleAll() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.NumberAttribute{
				Required: true,
				Validators: []validator.Number{
					// Validate this Number value must either be:
					//  - 1.0
					//  - 2.0, but not 3.0
					numbervalidator.Any(
						numbervalidator.OneOf(big.NewFloat(1.0)),
						numbervalidator.All(
							numbervalidator.OneOf(big.NewFloat(2.0)),
							numbervalidator.NoneOf(big.NewFloat(3.0)),
						),
					),
				},
			},
		},
	}
}
