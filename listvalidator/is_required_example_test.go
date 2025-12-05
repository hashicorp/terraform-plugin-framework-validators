// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package listvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleIsRequired() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Blocks: map[string]schema.Block{
			"example_block": schema.ListNestedBlock{
				Validators: []validator.List{
					// Validate this block has a value (not null).
					listvalidator.IsRequired(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"example_string_attribute": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}
