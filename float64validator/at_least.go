package float64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// floatAtLeastValidator validates that an float Attribute's value is at least a certain value.
type floatAtLeastValidator struct {
	tfsdk.AttributeValidator

	min float64
}

// Description describes the validation in plain text formatting.
func (validator floatAtLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at least %f", validator.min)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator floatAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator floatAtLeastValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	f, ok := validateFloat(ctx, request, response)
	if !ok {
		return
	}

	if f < validator.min {
		response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			request.AttributePath,
			"Invalid value",
			fmt.Sprintf("expected value to be at least %f, got %f", validator.min, f),
		))

		return
	}
}

// FloatAtLeast returns a new float value at least validator.
func FloatAtLeast(min float64) tfsdk.AttributeValidator {
	return floatAtLeastValidator{
		min: min,
	}
}
