// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package setvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
)

func ExamplePreferWriteOnlyAttribute() {
	// Used within a Schema method of a Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Set{
					// Throws a warning diagnostic encouraging practitioners to use
					// write_only_attr if example_attr has a value.
					setvalidator.PreferWriteOnlyAttribute(
						path.MatchRoot("write_only_attr"),
					),
				},
			},
			"write_only_attr": schema.SetAttribute{
				WriteOnly: true,
				Optional:  true,
			},
		},
	}
}
