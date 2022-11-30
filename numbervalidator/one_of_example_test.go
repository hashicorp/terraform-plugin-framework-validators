package numbervalidator_test

import (
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleOneOf() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.NumberAttribute{
				Required: true,
				Validators: []validator.Number{
					// Validate number value must be 1.2, 2.4, or 4.8
					numbervalidator.OneOf(
						[]*big.Float{
							big.NewFloat(1.2),
							big.NewFloat(2.4),
							big.NewFloat(4.8),
						}...,
					),
				},
			},
		},
	}
}
