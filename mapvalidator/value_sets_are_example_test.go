package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueSetsAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.MapAttribute{
				// This Map has values of Sets of Strings.
				// Roughly equivalent to map[string][]string.
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				Required: true,
				Validators: []validator.Map{
					// Validate this Map must contain Set elements
					// which have at least 1 String element.
					mapvalidator.ValueSetsAre(setvalidator.SizeAtLeast(1)),
				},
			},
		},
	}
}
