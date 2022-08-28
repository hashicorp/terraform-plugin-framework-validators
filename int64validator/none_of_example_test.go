package int64validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleNoneOf() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type:     types.Int64Type,
				Validators: []tfsdk.AttributeValidator{
					// Validate integer value must not be 12, 24, or 48
					int64validator.NoneOf([]int64{12, 24, 48}...),
				},
			},
		},
	}
}
