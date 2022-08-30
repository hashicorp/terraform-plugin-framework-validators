package float64validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleBetween() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type:     types.Float64Type,
				Validators: []tfsdk.AttributeValidator{
					// Validate floating point value must be at least 0.0 and at most 1.0
					float64validator.Between(0.0, 1.0),
				},
			},
		},
	}
}
