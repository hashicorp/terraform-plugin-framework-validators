package metavalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = allValidator{}

// allValidator validates that value validates against all the value validators and is for use in
// conjunction with the anyValidator, as the default behaviour is to validate all at the top-level.
type allValidator struct {
	valueValidators []tfsdk.AttributeValidator
}

// Description describes the validation in plain text formatting.
func (v allValidator) Description(ctx context.Context) string {
	var descriptions []string
	for _, validator := range v.valueValidators {
		descriptions = append(descriptions, validator.Description(ctx))
	}

	return fmt.Sprintf("Value must satisfy all of the validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v allValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
// If the number of iterations (i.e., k + 1) is greater than the number of diagnostics in the response then
// at least one of the validations has passed and, we return without any diagnostics.
func (v allValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	for _, validator := range v.valueValidators {
		validator.Validate(ctx, req, resp)
	}
}

// All returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Validates against all the value validators.
func All(valueValidators ...tfsdk.AttributeValidator) tfsdk.AttributeValidator {
	return allValidator{
		valueValidators: valueValidators,
	}
}
