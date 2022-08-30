package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValuesAre() {
	// Used within a GetSchema method of a DataSource, Provider, or Resource
	_ = tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"example_attr": {
				Required: true,
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Validators: []tfsdk.AttributeValidator{
					// Validate this map must contain string values which are at least 3 characters.
					mapvalidator.ValuesAre(stringvalidator.LengthAtLeast(3)),
				},
			},
		},
	}
}
