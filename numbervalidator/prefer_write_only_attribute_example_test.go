// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package numbervalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/numbervalidator"
)

func ExamplePreferWriteOnlyAttribute() {
	// Used within a Schema method of a Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.NumberAttribute{
				Optional: true,
				Validators: []validator.Number{
					// Throws a warning diagnostic encouraging practitioners to use
					// write_only_attr if example_attr has a value.
					numbervalidator.PreferWriteOnlyAttribute(
						path.MatchRoot("write_only_attr"),
					),
				},
			},
			"write_only_attr": schema.NumberAttribute{
				WriteOnly: true,
				Optional:  true,
			},
		},
	}
}
