package metavalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = anyWithAllWarningsValidator{}

// anyWithAllWarningsValidator validates that value validates against at least one of the value validators.
type anyWithAllWarningsValidator struct {
	valueValidators []tfsdk.AttributeValidator
}

// Description describes the validation in plain text formatting.
func (v anyWithAllWarningsValidator) Description(ctx context.Context) string {
	var descriptions []string
	for _, validator := range v.valueValidators {
		descriptions = append(descriptions, validator.Description(ctx))
	}

	return fmt.Sprintf("Value must satisfy at least one of the validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v anyWithAllWarningsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
// If the diagnostics returned from the validator do not contain an error we return all accumulated warning
// diagnostics from the passing validator(s) and  any failing validator(s).
func (v anyWithAllWarningsValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	anyValid := false

	for _, validator := range v.valueValidators {
		validatorResp := &tfsdk.ValidateAttributeResponse{
			Diagnostics: diag.Diagnostics{},
		}

		validator.Validate(ctx, req, validatorResp)

		if !validatorResp.Diagnostics.HasError() {
			anyValid = true
		}

		resp.Diagnostics.Append(validatorResp.Diagnostics...)
	}

	if anyValid {
		diagWarnings := diag.Diagnostics{}

		for _, d := range resp.Diagnostics {
			if d.Severity() == diag.SeverityWarning {
				diagWarnings.Append(d)
			}
		}

		resp.Diagnostics = diagWarnings
		return
	}
}

// AnyWithAllWarnings returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Validates against at least one of the value validators.
//	   - Returns all warnings for all passing and failing validators when at least one of the validators passes.
func AnyWithAllWarnings(valueValidators ...tfsdk.AttributeValidator) tfsdk.AttributeValidator {
	return anyWithAllWarningsValidator{
		valueValidators: valueValidators,
	}
}
