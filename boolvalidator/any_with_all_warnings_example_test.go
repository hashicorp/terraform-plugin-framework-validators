package boolvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
)

func ExampleAnyWithAllWarnings() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.BoolAttribute{
				Required: true,
				Validators: []validator.Bool{
					// Validate that this attribute must either:
					//  - be set with other_attrA or
					//  - be set with other_attrB
					boolvalidator.AnyWithAllWarnings(
						boolvalidator.AlsoRequires(path.Expressions{
							path.MatchRoot("other_attrA"),
						}...),
						boolvalidator.AlsoRequires(path.Expressions{
							path.MatchRoot("other_attrB"),
						}...),
					),
				},
			},
			"other_attrA": schema.StringAttribute{
				Optional: true,
			},
			"other_attrB": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}
