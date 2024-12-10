// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleNoNullValues() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Map{
					// Validate this map must contain no null values.
					mapvalidator.NoNullValues(),
				},
			},
		},
	}
}

func ExampleNoNullValues_function() {
	_ = function.Definition{
		Parameters: []function.Parameter{
			function.MapParameter{
				Name: "example_param",
				Validators: []function.MapParameterValidator{
					// Validate this map must contain no null values.
					mapvalidator.NoNullValues(),
				},
			},
		},
	}
}
