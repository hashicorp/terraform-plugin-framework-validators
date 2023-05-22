// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package setvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleAny() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.SetAttribute{
				Required: true,
				Validators: []validator.Set{
					// Validate this Set value must either be:
					//  - Between 1 and 2 elements
					//  - At least 4 elements
					setvalidator.Any(
						setvalidator.SizeBetween(1, 2),
						setvalidator.SizeAtLeast(4),
					),
				},
			},
		},
	}
}
