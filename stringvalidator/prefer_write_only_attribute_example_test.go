// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package stringvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func ExamplePreferWriteOnlyAttribute() {
	// Used within a Schema method of a Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					// Throws a warning diagnostic encouraging practitioners to use
					// write_only_attr if example_attr has a value.
					stringvalidator.PreferWriteOnlyAttribute(
						path.MatchRoot("write_only_attr"),
					),
				},
			},
			"write_only_attr": schema.StringAttribute{
				WriteOnly: true,
				Optional:  true,
			},
		},
	}
}
