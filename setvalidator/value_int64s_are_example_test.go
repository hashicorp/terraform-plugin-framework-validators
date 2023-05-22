// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package setvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleValueInt64sAre() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.SetAttribute{
				ElementType: types.Int64Type,
				Required:    true,
				Validators: []validator.Set{
					// Validate this Set must contain Int64 values which are at least 1.
					setvalidator.ValueInt64sAre(int64validator.AtLeast(1)),
				},
			},
		},
	}
}
