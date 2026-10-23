// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package resourcevalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// PreferWriteOnlyAttributeIfFunc is a conditional function used in the
// PreferWriteOnlyAttributeIf validator to determine whether a warning should be returned.
type PreferWriteOnlyAttributeIfFunc func(context.Context, resource.ValidateConfigRequest, *PreferWriteOnlyAttributeIfFuncResponse)

// PreferWriteOnlyAttributeIfFuncResponse is the response type for a
// PreferWriteOnlyAttributeIfFunc.
type PreferWriteOnlyAttributeIfFuncResponse struct {
	// Diagnostics report errors or warnings related to the conditional logic.
	Diagnostics diag.Diagnostics

	// Warn enables the PreferWriteOnlyAttribute validation.
	Warn bool
}

// PreferWriteOnlyAttributeIf returns a warning if the given function enables it
// and the PreferWriteOnlyAttribute validation conditions are met.
func PreferWriteOnlyAttributeIf(oldAttribute path.Expression, writeOnlyAttribute path.Expression, f PreferWriteOnlyAttributeIfFunc) resource.ConfigValidator {
	return preferWriteOnlyAttributeIfValidator{
		ifFunc:    f,
		validator: PreferWriteOnlyAttribute(oldAttribute, writeOnlyAttribute),
	}
}

var _ resource.ConfigValidator = preferWriteOnlyAttributeIfValidator{}

type preferWriteOnlyAttributeIfValidator struct {
	ifFunc    PreferWriteOnlyAttributeIfFunc
	validator resource.ConfigValidator
}

func (v preferWriteOnlyAttributeIfValidator) Description(ctx context.Context) string {
	return v.validator.Description(ctx)
}

func (v preferWriteOnlyAttributeIfValidator) MarkdownDescription(ctx context.Context) string {
	return v.validator.MarkdownDescription(ctx)
}

func (v preferWriteOnlyAttributeIfValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	ifResp := &PreferWriteOnlyAttributeIfFuncResponse{}
	v.ifFunc(ctx, req, ifResp)
	resp.Diagnostics.Append(ifResp.Diagnostics...)
	if ifResp.Warn {
		v.validator.ValidateResource(ctx, req, resp)
	}
}
