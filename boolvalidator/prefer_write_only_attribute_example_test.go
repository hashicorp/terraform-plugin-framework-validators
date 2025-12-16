// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package boolvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
)

func ExamplePreferWriteOnlyAttribute() {
	// Used within a Schema method of a Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.BoolAttribute{
				Optional: true,
				Validators: []validator.Bool{
					// Throws a warning diagnostic encouraging practitioners to use
					// write_only_attr if example_attr has a value.
					boolvalidator.PreferWriteOnlyAttribute(
						path.MatchRoot("write_only_attr"),
					),
				},
			},
			"write_only_attr": schema.BoolAttribute{
				WriteOnly: true,
				Optional:  true,
			},
		},
	}
}
