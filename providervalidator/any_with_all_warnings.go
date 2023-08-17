package providervalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// AnyWithAllWarnings returns a validator which ensures that any configured
// attribute value passes at least one of the given validators. This validator
// returns all warnings, including failed validators.
//
// Use Any() to return warnings only from the passing validator.
func AnyWithAllWarnings(validators ...provider.ConfigValidator) provider.ConfigValidator {
	return anyWithAllWarningsValidator{
		validators: validators,
	}
}

var _ provider.ConfigValidator = anyWithAllWarningsValidator{}

// anyWithAllWarningsValidator implements the validator.
type anyWithAllWarningsValidator struct {
	validators []provider.ConfigValidator
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

// ValidateProvider performs the validation.
func (v anyWithAllWarningsValidator) ValidateProvider(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	anyValid := false

	for _, subValidator := range v.validators {
		validateResp := &provider.ValidateConfigResponse{}

		subValidator.ValidateProvider(ctx, req, validateResp)

		if !validateResp.Diagnostics.HasError() {
			anyValid = true
		}

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}

	if anyValid {
		resp.Diagnostics = resp.Diagnostics.Warnings()
	}
}
