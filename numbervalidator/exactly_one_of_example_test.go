// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package numbervalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleExactlyOneOf() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.NumberAttribute{
				Optional: true,
				Validators: []validator.Number{
					// Validate only this attribute or other_attr is configured.
					numbervalidator.ExactlyOneOf(path.Expressions{
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
