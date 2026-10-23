// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package stringvalidator_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
					// example_attr_wo if example_attr has a value.
					stringvalidator.PreferWriteOnlyAttribute(
						path.MatchRoot("example_attr_wo"),
					),
				},
			},
			"example_attr_wo": schema.StringAttribute{
				WriteOnly: true,
				Optional:  true,
			},
		},
	}
}

func ExamplePreferWriteOnlyAttributeIf() {
	// Used within a Schema method of a Resource
	_ = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"example_attr": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					// Throws a warning diagnostic encouraging practitioners to use
					// example_attr_wo if example_attr has a value and sensitive is true.
					stringvalidator.PreferWriteOnlyAttributeIf(
						path.MatchRoot("example_attr_wo"),
						func(ctx context.Context, req validator.StringRequest, resp *stringvalidator.PreferWriteOnlyAttributeIfFuncResponse) {
							var sensitive types.Bool
							diags := req.Config.GetAttribute(ctx, path.Root("sensitive"), &sensitive)
							resp.Diagnostics.Append(diags...)
							resp.Warn = !diags.HasError() && !sensitive.IsNull() && !sensitive.IsUnknown() && sensitive.ValueBool()
						},
					),
				},
			},
			"example_attr_wo": schema.StringAttribute{
				WriteOnly: true,
				Optional:  true,
			},
			"sensitive": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}
