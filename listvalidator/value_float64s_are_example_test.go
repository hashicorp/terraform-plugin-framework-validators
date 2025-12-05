// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package listvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueFloat64sAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.ListAttribute{
				ElementType: types.Float64Type,
				Required:    true,
				Validators: []validator.List{
					// Validate this List must contain Float64 values which are at least 1.2.
					listvalidator.ValueFloat64sAre(float64validator.AtLeast(1.2)),
				},
			},
		},
	}
}
