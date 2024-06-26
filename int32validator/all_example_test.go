// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int32validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
)

func ExampleAll() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.Int32Attribute{
				Required: true,
				Validators: []validator.Int32{
					// Validate this Int32 value must either be:
					//  - 1.0
					//  - At least 2.0, but not 3.0
					int32validator.Any(
						int32validator.OneOf(1.0),
						int32validator.All(
							int32validator.AtLeast(2.0),
							int32validator.NoneOf(3.0),
						),
					),
				},
			},
		},
	}
}
