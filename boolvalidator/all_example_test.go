// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boolvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
)

func ExampleAll() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.BoolAttribute{
				Required: true,
				Validators: []validator.Bool{
					// This attribute must satisfy either All validator.
					boolvalidator.Any(
						boolvalidator.All( /* ... */ ),
						boolvalidator.All( /* ... */ ),
					),
				},
			},
		},
	}
}
