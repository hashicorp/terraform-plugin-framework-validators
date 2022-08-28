package setvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleSizeAtMost() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Validators: []tfsdk.AttributeValidator{
					// Validate this set must contain at most 2 elements.
					setvalidator.SizeAtMost(2),
				},
			},
		},
	}
}
