package stringvalidator

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = regexMatchesValidator{}

// regexMatchesValidator validates that a string Attribute's value matches the specified regular expression.
type regexMatchesValidator struct {
	regexp  *regexp.Regexp
	message string
}

// Description describes the validation in plain text formatting.
func (validator regexMatchesValidator) Description(_ context.Context) string {
	if validator.message != "" {
		return validator.message
	}
	return fmt.Sprintf("value must match regular expression '%s'", validator.regexp)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator regexMatchesValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator regexMatchesValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	s, ok := validateString(ctx, request, response)

	if !ok {
		return
	}

	if ok := validator.regexp.MatchString(s); !ok {
		response.Diagnostics.Append(validatordiag.AttributeValueMatchesDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			s,
		))
	}
}

// RegexMatches returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Is a string.
//     - Matches the given regular expression.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func RegexMatches(regexp *regexp.Regexp, message string) tfsdk.AttributeValidator {
	return regexMatchesValidator{
		regexp:  regexp,
		message: message,
	}
}
