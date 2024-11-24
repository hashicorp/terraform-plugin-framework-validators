// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynamicvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Any returns a validator which ensures that any configured attribute value
// passes at least one of the given validators.
//
// To prevent practitioner confusion should non-passing validators have
// conflicting logic, only warnings from the passing validator are returned.
// Use AnyWithAllWarnings() to return warnings from non-passing validators
// as well.
func Any(validators ...validator.Dynamic) validator.Dynamic {
	return anyValidator{
		validators: validators,
	}
}

var _ validator.Dynamic = anyValidator{}

// anyValidator implements the validator.
type anyValidator struct {
	validators []validator.Dynamic
}

// Description describes the validation in plain text formatting.
func (v anyValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, subValidator := range v.validators {
		descriptions = append(descriptions, subValidator.Description(ctx))
	}

	return fmt.Sprintf("Value must satisfy at least one of the validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v anyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateDynamic performs the validation.
func (v anyValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	for _, subValidator := range v.validators {
		validateResp := &validator.DynamicResponse{}

		subValidator.ValidateDynamic(ctx, req, validateResp)

		if !validateResp.Diagnostics.HasError() {
			resp.Diagnostics = validateResp.Diagnostics

			return
		}

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}
