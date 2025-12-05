// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package setvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleAlsoRequires() {
	// Used within a Schema method of a DataSource, Provider, or Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Set{
					// Validate this attribute must be configured with other_attr.
					setvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("other_attr"),
					}...),
				},
			},
			"other_attr": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}
