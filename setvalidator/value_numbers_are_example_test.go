package setvalidator_test

import (
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueNumbersAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.SetAttribute{
				ElementType: types.NumberType,
				Required:    true,
				Validators: []validator.Set{
					// Validate this Set must contain Number values which are 1.2 or 2.4.
					setvalidator.ValueNumbersAre(
						numbervalidator.OneOf(
							big.NewFloat(1.2),
							big.NewFloat(2.4),
						),
					),
				},
			},
		},
	}
}
