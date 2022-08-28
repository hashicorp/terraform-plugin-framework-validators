package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleSizeAtMost() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Validators: []tfsdk.AttributeValidator{
					// Validate this map must contain at most 2 elements.
					mapvalidator.SizeAtMost(2),
				},
			},
		},
	}
}
