// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int32validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
)

func ExampleBetween() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.Int32Attribute{
				Required: true,
				Validators: []validator.Int32{
					// Validate integer value must be at least 10 and at most 100
					int32validator.Between(10, 100),
				},
			},
		},
	}
}
