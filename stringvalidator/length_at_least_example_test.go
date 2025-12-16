// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package stringvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleLengthAtLeast() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					// Validate string value length must be at least 3 characters.
					stringvalidator.LengthAtLeast(3),
				},
			},
		},
	}
}

func ExampleLengthAtLeast_function() {
	_ = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				Name: "example_param",
				Validators: []function.StringParameterValidator{
					// Validate string value length must be at least 3 characters.
					stringvalidator.LengthAtLeast(3),
				},
			},
		},
	}
}
