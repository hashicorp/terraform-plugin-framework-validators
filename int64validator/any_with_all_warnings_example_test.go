package int64validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleAnyWithAllWarnings() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					// Validate this Int64 value must either be:
					//  - 1.0
					//  - At least 2.0
					int64validator.AnyWithAllWarnings(
						int64validator.OneOf(1.0),
						int64validator.AtLeast(2.0),
					),
				},
			},
		},
	}
}
