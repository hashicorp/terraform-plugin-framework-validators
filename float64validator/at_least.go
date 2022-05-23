package float64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = atLeastValidator{}

// atLeastValidator validates that an float Attribute's value is at least a certain value.
type atLeastValidator struct {
	min float64
}

// Description describes the validation in plain text formatting.
func (validator atLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at least %f", validator.min)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator atLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator atLeastValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	f, ok := validateFloat(ctx, request, response)
	if !ok {
		return
	}

	if f < validator.min {
		response.Diagnostics.Append(validatordiag.AttributeValueDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			fmt.Sprintf("%f", f),
		))

		return
	}
}

// AtLeast returns a new float value at least validator.
func AtLeast(min float64) tfsdk.AttributeValidator {
	return atLeastValidator{
		min: min,
	}
}
