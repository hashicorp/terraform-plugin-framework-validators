// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package setvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleAll() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Set{
					// Validate this Set value must either be:
					//  - More than 5 elements
					//  - At least 2 elements, but not more than 3 elements
					setvalidator.Any(
						setvalidator.SizeAtLeast(5),
						setvalidator.All(
							setvalidator.SizeAtLeast(2),
							setvalidator.SizeAtMost(3),
						),
					),
				},
			},
		},
	}
}
