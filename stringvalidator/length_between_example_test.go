// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleLengthBetween() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					// Validate string value length must be at least 3 and at most 256 characters.
					stringvalidator.LengthBetween(3, 256),
				},
			},
		},
	}
}

func ExampleLengthBetween_function() {
	_ = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				Name: "example_param",
				Validators: []function.StringParameterValidator{
					// Validate string value length must be at least 3 and at most 256 characters.
					stringvalidator.LengthBetween(3, 256),
				},
			},
		},
	}
}
