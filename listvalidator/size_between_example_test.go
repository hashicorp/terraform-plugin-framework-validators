// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleSizeBetween() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					// Validate this list must contain at least 2 and at most 4 elements.
					listvalidator.SizeBetween(2, 4),
				},
			},
		},
	}
}

func ExampleSizeBetween_function() {
	_ = function.Definition{
		Parameters: []function.Parameter{
			function.ListParameter{
				Name: "example_param",
				Validators: []function.ListParameterValidator{
					// Validate this list must contain at least 2 and at most 4 elements.
					listvalidator.SizeBetween(2, 4),
				},
			},
		},
	}
}
