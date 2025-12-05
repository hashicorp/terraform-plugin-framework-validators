// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package int32validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
)

func ExampleOneOf() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.Int32Attribute{
				Required: true,
				Validators: []validator.Int32{
					// Validate integer value must be 12, 24, or 48
					int32validator.OneOf([]int32{12, 24, 48}...),
				},
			},
		},
	}
}

func ExampleOneOf_function() {
	_ = function.Definition{
		Parameters: []function.Parameter{
			function.Int32Parameter{
				Name: "example_param",
				Validators: []function.Int32ParameterValidator{
					// Validate integer value must be 12, 24, or 48
					int32validator.OneOf([]int32{12, 24, 48}...),
				},
			},
		},
	}
}
