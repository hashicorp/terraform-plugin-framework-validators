// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueStringsAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Map{
					// Validate this Map must contain string values which are at least 3 characters.
					mapvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(3)),
				},
			},
		},
	}
}
