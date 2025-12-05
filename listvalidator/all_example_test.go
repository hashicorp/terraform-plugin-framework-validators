// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package listvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleAll() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					// Validate this List value must either be:
					//  - More than 5 elements
					//  - At least 2 elements, but not more than 3 elements
					listvalidator.Any(
						listvalidator.SizeAtLeast(5),
						listvalidator.All(
							listvalidator.SizeAtLeast(2),
							listvalidator.SizeAtMost(3),
						),
					),
				},
			},
		},
	}
}
