package listvalidator_test

import (
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueNumbersAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.ListAttribute{
				ElementType: types.NumberType,
				Required:    true,
				Validators: []validator.List{
					// Validate this List must contain Number values which are 1.2 or 2.4.
					listvalidator.ValueNumbersAre(
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
