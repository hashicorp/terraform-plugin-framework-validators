package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueMapsAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.MapAttribute{
				// This Map has values of Maps of Strings.
				// Roughly equivalent to map[string]map[string]string.
				ElementType: types.MapType{
					ElemType: types.StringType,
				},
				Required: true,
				Validators: []validator.Map{
					// Validate this Map must contain Map elements
					// which have at least 1 element.
					mapvalidator.ValueMapsAre(mapvalidator.SizeAtLeast(1)),
				},
			},
		},
	}
}
