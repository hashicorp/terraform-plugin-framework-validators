// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynamicvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// AnyWithAllWarnings returns a validator which ensures that any configured
// attribute value passes at least one of the given validators. This validator
// returns all warnings, including failed validators.
//
// Use Any() to return warnings only from the passing validator.
func AnyWithAllWarnings(validators ...validator.Dynamic) validator.Dynamic {
	return anyWithAllWarningsValidator{
		validators: validators,
	}
}

var _ validator.Dynamic = anyWithAllWarningsValidator{}

// anyWithAllWarningsValidator implements the validator.
type anyWithAllWarningsValidator struct {
	validators []validator.Dynamic
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

// ValidateDynamic performs the validation.
func (v anyWithAllWarningsValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	anyValid := false

	for _, subValidator := range v.validators {
		validateResp := &validator.DynamicResponse{}

		subValidator.ValidateDynamic(ctx, req, validateResp)

		if !validateResp.Diagnostics.HasError() {
			anyValid = true
		}

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}

	if anyValid {
		resp.Diagnostics = resp.Diagnostics.Warnings()
	}
}
