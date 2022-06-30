package metavalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = allValidator{}

// allValidator validates that value validates against all the value validators and is for use in
// conjunction with the anyValidator or anyWithAllWarningsValidator, as the default behaviour is
// to validate all at the top-level.
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
func (v allValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	for _, validator := range v.valueValidators {
		validator.Validate(ctx, req, resp)
	}
}

// All returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Validates against all the value validators.
//
// Use of All is only necessary when used in conjunction with Any or AnyWithAllWarnings
// as the []tfsdk.AttributeValidator field automatically applies a logical AND.
func All(valueValidators ...tfsdk.AttributeValidator) tfsdk.AttributeValidator {
	return allValidator{
		valueValidators: valueValidators,
	}
}
