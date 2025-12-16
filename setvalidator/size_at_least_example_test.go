// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package setvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleSizeAtLeast() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Set{
					// Validate this set must contain at least 2 elements.
					setvalidator.SizeAtLeast(2),
				},
			},
		},
	}
}

func ExampleSizeAtLeast_function() {
	_ = function.Definition{
		Parameters: []function.Parameter{
			function.SetParameter{
				Name: "example_param",
				Validators: []function.SetParameterValidator{
					// Validate this set must contain at least 2 elements.
					setvalidator.SizeAtLeast(2),
				},
			},
		},
	}
}
