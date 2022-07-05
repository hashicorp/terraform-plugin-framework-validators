package metavalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = anyValidator{}

// anyValidator validates that value validates against at least one of the value validators.
type anyValidator struct {
	valueValidators []tfsdk.AttributeValidator
}

// Description describes the validation in plain text formatting.
func (v anyValidator) Description(ctx context.Context) string {
	var descriptions []string
	for _, validator := range v.valueValidators {
		descriptions = append(descriptions, validator.Description(ctx))
	}

	return fmt.Sprintf("Value must satisfy at least one of the validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v anyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
// The validator will pass if it encounters a value validator that returns no errors and will then return any warnings
// from the passing validator. Using All validator as value validators will pass if all the validators supplied in an
// All validator pass.
func (v anyValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	for _, validator := range v.valueValidators {
		validatorResp := &tfsdk.ValidateAttributeResponse{
			Diagnostics: diag.Diagnostics{},
		}

		validator.Validate(ctx, req, validatorResp)

		if !validatorResp.Diagnostics.HasError() {
			resp.Diagnostics = validatorResp.Diagnostics
			return
		}

		resp.Diagnostics.Append(validatorResp.Diagnostics...)
	}
}

// Any returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Validates against at least one of the value validators.
//
// To prevent practitioner confusion should non-passing validators have
// conflicting logic, only warnings from the passing validator are returned.
// Use AnyWithAllWarnings() to return warnings from non-passing validators
// as well.
func Any(valueValidators ...tfsdk.AttributeValidator) tfsdk.AttributeValidator {
	return anyValidator{
		valueValidators: valueValidators,
	}
}
