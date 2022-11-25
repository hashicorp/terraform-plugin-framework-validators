package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueFloat64sAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.MapAttribute{
				ElementType: types.Float64Type,
				Required:    true,
				Validators: []validator.Map{
					// Validate this Map must contain Float64 values which are at least 1.2.
					mapvalidator.ValueFloat64sAre(float64validator.AtLeast(1.2)),
				},
			},
		},
	}
}
