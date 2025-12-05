// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package stringvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleUTF8LengthBetween() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					// Validate UTF-8 character count must be at least 3 characters
					// and at most 255 characters.
					stringvalidator.UTF8LengthBetween(3, 255),
				},
			},
		},
	}
}

func ExampleUTF8LengthBetween_function() {
	_ = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				Name: "example_param",
				Validators: []function.StringParameterValidator{
					// Validate UTF-8 character count must be at least 3 characters
					// and at most 255 characters.
					stringvalidator.UTF8LengthBetween(3, 255),
				},
			},
		},
	}
}
