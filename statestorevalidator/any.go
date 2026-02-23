// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package statestorevalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

// Any returns a validator which ensures that any configured attribute value
// passes at least one of the given validators.
//
// To prevent practitioner confusion should non-passing validators have
// conflicting logic, only warnings from the passing validator are returned.
// Use AnyWithAllWarnings() to return warnings from non-passing validators
// as well.
func Any(validators ...statestore.ConfigValidator) statestore.ConfigValidator {
	return anyValidator{
		validators: validators,
	}
}

var _ statestore.ConfigValidator = anyValidator{}

// anyValidator implements the validator.
type anyValidator struct {
	validators []statestore.ConfigValidator
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

// ValidateStateStore performs the validation.
func (v anyValidator) ValidateStateStore(ctx context.Context, req statestore.ValidateConfigRequest, resp *statestore.ValidateConfigResponse) {
	for _, subValidator := range v.validators {
		validateResp := &statestore.ValidateConfigResponse{}

		subValidator.ValidateStateStore(ctx, req, validateResp)

		if !validateResp.Diagnostics.HasError() {
			resp.Diagnostics = validateResp.Diagnostics

			return
		}

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}
