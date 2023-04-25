package int64validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
)

func ExampleEqualToProductOf() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					// Validate this integer value must be equal to the
					// product of integer values other_attr1 and other_attr2.
					int64validator.EqualToProductOf(path.Expressions{
						path.MatchRoot("other_attr1"),
						path.MatchRoot("other_attr2"),
					}...),
				},
			},
			"other_attr1": schema.Int64Attribute{
				Required: true,
			},
			"other_attr2": schema.Int64Attribute{
				Required: true,
			},
		},
	}
}
