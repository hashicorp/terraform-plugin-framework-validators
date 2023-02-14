package int64validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleDifferentFrom() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					// Validate this attribute must be configured with other_attr.
					int64validator.DifferentFrom(path.Expressions{
						path.MatchRoot("other_attr"),
					}...),
				},
			},
			"other_attr": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}
