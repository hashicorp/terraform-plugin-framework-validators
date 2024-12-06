// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynamicvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/dynamicvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleAll() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.DynamicAttribute{
				Required: true,
				Validators: []validator.Dynamic{
					// Validate this Dynamic value must either be:
					//  - "one"
					dynamicvalidator.Any(
						dynamicvalidator.Any( /* ... */ ),
						dynamicvalidator.All( /* ... */ ),
					),
				},
			},
		},
	}
}
