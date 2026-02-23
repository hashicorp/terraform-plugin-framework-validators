// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package statestorevalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

// All returns a validator which ensures that any configured attribute value
// validates against all the given validators.
//
// Use of All is only necessary when used in conjunction with Any or AnyWithAllWarnings
// as the Validators field automatically applies a logical AND.
func All(validators ...statestore.ConfigValidator) statestore.ConfigValidator {
	return allValidator{
		validators: validators,
	}
}

var _ statestore.ConfigValidator = allValidator{}

// allValidator implements the validator.
type allValidator struct {
	validators []statestore.ConfigValidator
}

// Description describes the validation in plain text formatting.
func (v allValidator) Description(ctx context.Context) string {
	var descriptions []string

	for _, subValidator := range v.validators {
		descriptions = append(descriptions, subValidator.Description(ctx))
	}

	return fmt.Sprintf("Value must satisfy all of the validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v allValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateStateStore performs the validation.
func (v allValidator) ValidateStateStore(ctx context.Context, req statestore.ValidateConfigRequest, resp *statestore.ValidateConfigResponse) {
	for _, subValidator := range v.validators {
		validateResp := &statestore.ValidateConfigResponse{}

		subValidator.ValidateStateStore(ctx, req, validateResp)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}
