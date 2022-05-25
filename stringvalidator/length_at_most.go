package stringvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.AttributeValidator = lengthAtMostValidator{}

// lengthAtMostValidator validates that a string Attribute's length is at most a certain value.
type lengthAtMostValidator struct {
	maxLength int
}

// Description describes the validation in plain text formatting.
func (validator lengthAtMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("string length must be at most %d", validator.maxLength)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator lengthAtMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator lengthAtMostValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	s, ok := validateString(ctx, request, response)

	if !ok {
		return
	}

	if l := len(s); l > validator.maxLength {
		response.Diagnostics.Append(validatordiag.AttributeValueLengthDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			fmt.Sprintf("%d", l),
		))

		return
	}
}

// LengthAtMost returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Is a string.
//     - Is of length exclusively less than the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func LengthAtMost(maxLength int) tfsdk.AttributeValidator {
	if maxLength < 0 {
		return nil
	}

	return lengthAtMostValidator{
		maxLength: maxLength,
	}
}

// validateString ensures that the request contains a String value.
func validateString(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) (string, bool) {
	var s types.String

	diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &s)

	if diags.HasError() {
		response.Diagnostics = append(response.Diagnostics, diags...)

		return "", false
	}

	if s.Unknown || s.Null {
		return "", false
	}

	return s.Value, true
}
