package float64validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleAtLeast() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type:     types.Float64Type,
				Validators: []tfsdk.AttributeValidator{
					// Validate floating point value must be at least 42.42
					float64validator.AtLeast(42.42),
				},
			},
		},
	}
}
