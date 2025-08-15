// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actionvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
)

// AnyWithAllWarnings returns a validator which ensures that any configured
// attribute value passes at least one of the given validators. This validator
// returns all warnings, including failed validators.
//
// Use Any() to return warnings only from the passing validator.
func AnyWithAllWarnings(validators ...action.ConfigValidator) action.ConfigValidator {
	return anyWithAllWarningsValidator{
		validators: validators,
	}
}

var _ action.ConfigValidator = anyWithAllWarningsValidator{}

// anyWithAllWarningsValidator implements the validator.
type anyWithAllWarningsValidator struct {
	validators []action.ConfigValidator
}

// Description describes the validation in plain text formatting.
func (v anyWithAllWarningsValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, subValidator := range v.validators {
		descriptions = append(descriptions, subValidator.Description(ctx))
	}

	return fmt.Sprintf("Value must satisfy at least one of the validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v anyWithAllWarningsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateAction performs the validation.
func (v anyWithAllWarningsValidator) ValidateAction(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	anyValid := false

	for _, subValidator := range v.validators {
		validateResp := &action.ValidateConfigResponse{}

		subValidator.ValidateAction(ctx, req, validateResp)

		if !validateResp.Diagnostics.HasError() {
			anyValid = true
		}

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}

	if anyValid {
		resp.Diagnostics = resp.Diagnostics.Warnings()
	}
}
