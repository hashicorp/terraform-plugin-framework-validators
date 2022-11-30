package listvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueSetsAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.ListAttribute{
				// This List has values of Sets of Strings.
				// Roughly equivalent to [][]string.
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				Required: true,
				Validators: []validator.List{
					// Validate this List must contain Set elements
					// which have at least 1 String element.
					listvalidator.ValueSetsAre(setvalidator.SizeAtLeast(1)),
				},
			},
		},
	}
}
