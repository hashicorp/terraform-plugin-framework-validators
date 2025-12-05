// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package setvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
)

func ExampleValueInt32sAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.SetAttribute{
				ElementType: types.Int32Type,
				Required:    true,
				Validators: []validator.Set{
					// Validate this Set must contain Int32 values which are at least 1.
					setvalidator.ValueInt32sAre(int32validator.AtLeast(1)),
				},
			},
		},
	}
}
