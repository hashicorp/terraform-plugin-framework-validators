// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package float32validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/float32validator"
)

func ExampleAll() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.Float32Attribute{
				Required: true,
				Validators: []validator.Float32{
					// Validate this Float32 value must either be:
					//  - 1.0
					//  - At least 2.0, but not 3.0
					float32validator.Any(
						float32validator.OneOf(1.0),
						float32validator.All(
							float32validator.AtLeast(2.0),
							float32validator.NoneOf(3.0),
						),
					),
				},
			},
		},
	}
}
