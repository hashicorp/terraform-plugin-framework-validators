// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float64validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleAny() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.Float64Attribute{
				Required: true,
				Validators: []validator.Float64{
					// Validate this Float64 value must either be:
					//  - 1.0
					//  - At least 2.0
					float64validator.Any(
						float64validator.OneOf(1.0),
						float64validator.AtLeast(2.0),
					),
				},
			},
		},
	}
}
