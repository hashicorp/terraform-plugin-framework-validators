// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package stringvalidator_test

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ExampleRegexMatches() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					// Validate string value satisfies the regular expression for alphanumeric characters
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9]*$`),
						"must only contain only alphanumeric characters",
					),
				},
			},
		},
	}
}

func ExampleRegexMatches_function() {
	_ = function.Definition{
		Parameters: []function.Parameter{
			function.StringParameter{
				Name: "example_param",
				Validators: []function.StringParameterValidator{
					// Validate string value satisfies the regular expression for alphanumeric characters
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9]*$`),
						"must only contain only alphanumeric characters",
					),
				},
			},
		},
	}
}
