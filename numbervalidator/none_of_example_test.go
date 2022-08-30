package numbervalidator_test

import (
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleNoneOf() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type:     types.NumberType,
				Validators: []tfsdk.AttributeValidator{
					// Validate number value must not be 1.2, 2.4, or 4.8
					numbervalidator.NoneOf(
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
