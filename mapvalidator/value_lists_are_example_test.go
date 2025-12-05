// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueListsAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.MapAttribute{
				// This Map has values of Lists of Strings.
				// Roughly equivalent to map[string][]string.
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Required: true,
				Validators: []validator.Map{
					// Validate this Map must contain List elements
					// which have at least 1 String element.
					mapvalidator.ValueListsAre(listvalidator.SizeAtLeast(1)),
				},
			},
		},
	}
}
