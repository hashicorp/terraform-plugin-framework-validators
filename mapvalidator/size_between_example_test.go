// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package mapvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleSizeBetween() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Map{
					// Validate this map must contain at least 2 and at most 4 elements.
					mapvalidator.SizeBetween(2, 4),
				},
			},
		},
	}
}

func ExampleSizeBetween_function() {
	_ = function.Definition{
		Parameters: []function.Parameter{
			function.MapParameter{
				Name: "example_param",
				Validators: []function.MapParameterValidator{
					// Validate this map must contain at least 2 and at most 4 elements.
					mapvalidator.SizeBetween(2, 4),
				},
			},
		},
	}
}
