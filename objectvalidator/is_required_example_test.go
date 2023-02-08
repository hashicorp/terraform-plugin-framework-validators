package objectvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleIsRequired() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Blocks: map[string]schema.Block{
			"example_block": schema.SingleNestedBlock{
				Validators: []validator.Object{
					// Validate this block has a value (not null).
					objectvalidator.IsRequired(),
				},
				Attributes: map[string]schema.Attribute{
					"example_string_attribute": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}
