package float64validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleNoneOf() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type:     types.Float64Type,
				Validators: []tfsdk.AttributeValidator{
					// Validate floating point value must not be 1.2, 2.4, or 4.8
					float64validator.NoneOf([]float64{1.2, 2.4, 4.8}...),
				},
			},
		},
	}
}
