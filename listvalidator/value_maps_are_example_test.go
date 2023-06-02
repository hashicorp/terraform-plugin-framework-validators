// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueMapsAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.ListAttribute{
				// This List has values of Map of Strings.
				// Roughly equivalent to []map[string]string.
				ElementType: types.MapType{
					ElemType: types.StringType,
				},
				Required: true,
				Validators: []validator.List{
					// Validate this List must contain Map elements
					// which have at least 1 element.
					listvalidator.ValueMapsAre(mapvalidator.SizeAtLeast(1)),
				},
			},
		},
	}
}
