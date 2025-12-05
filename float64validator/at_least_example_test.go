// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package float64validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleAtLeast() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.Float64Attribute{
				Required: true,
				Validators: []validator.Float64{
					// Validate floating point value must be at least 42.42
					float64validator.AtLeast(42.42),
				},
			},
		},
	}
}

func ExampleAtLeast_function() {
	_ = function.Definition{
		Parameters: []function.Parameter{
			function.Float64Parameter{
				Name: "example_param",
				Validators: []function.Float64ParameterValidator{
					// Validate floating point value must be at least 42.42
					float64validator.AtLeast(42.42),
				},
			},
		},
	}
}
