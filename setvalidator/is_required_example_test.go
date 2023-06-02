// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package setvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleIsRequired() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Blocks: map[string]schema.Block{
			"example_block": schema.SetNestedBlock{
				Validators: []validator.Set{
					// Validate this block has a value (not null).
					setvalidator.IsRequired(),
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
