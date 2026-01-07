// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package numbervalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// All returns a validator which ensures that any configured attribute value
// attribute value validates against all the given validators.
//
// Use of All is only necessary when used in conjunction with Any or AnyWithAllWarnings
// as the Validators field automatically applies a logical AND.
func All(validators ...validator.Number) validator.Number {
	return allValidator{
		validators: validators,
	}
}

var _ validator.Number = allValidator{}

// allValidator implements the validator.
type allValidator struct {
	validators []validator.Number
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

// ValidateNumber performs the validation.
func (v allValidator) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	for _, subValidator := range v.validators {
		validateResp := &validator.NumberResponse{}

		subValidator.ValidateNumber(ctx, req, validateResp)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}
