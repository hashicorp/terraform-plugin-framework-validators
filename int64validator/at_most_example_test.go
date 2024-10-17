// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int64validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleAtMost() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					// Validate integer value must be at most 42
					int64validator.AtMost(42),
				},
			},
		},
	}
}

func ExampleAtMost_function() {
	_ = function.Definition{
		Parameters: []function.Parameter{
			function.Int64Parameter{
				Name: "example_param",
				Validators: []function.Int64ParameterValidator{
					// Validate integer value must be at most 42
					int64validator.AtMost(42),
				},
			},
		},
	}
}
