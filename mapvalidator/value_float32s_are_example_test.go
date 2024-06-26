// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/float32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
)

func ExampleValueFloat32sAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.MapAttribute{
				ElementType: types.Float32Type,
				Required:    true,
				Validators: []validator.Map{
					// Validate this Map must contain Float32 values which are at least 1.2.
					mapvalidator.ValueFloat32sAre(float32validator.AtLeast(1.2)),
				},
			},
		},
	}
}
