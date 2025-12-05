// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package listvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueStringsAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					// Validate this List must contain string values which are at least 3 characters.
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(3)),
				},
			},
		},
	}
}
