// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int32validator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
)

func ExampleAtMostSumOf() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.Int32Attribute{
				Required: true,
				Validators: []validator.Int32{
					// Validate this integer value must be at most the
					// summed integer values of other_attr1 and other_attr2.
					int32validator.AtMostSumOf(path.Expressions{
						path.MatchRoot("other_attr1"),
						path.MatchRoot("other_attr2"),
					}...),
				},
			},
			"other_attr1": schema.Int32Attribute{
				Required: true,
			},
			"other_attr2": schema.Int32Attribute{
				Required: true,
			},
		},
	}
}
