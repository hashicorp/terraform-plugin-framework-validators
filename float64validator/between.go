package float64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// betweenValidator validates that an float Attribute's value is in a range.
type betweenValidator struct {
	tfsdk.AttributeValidator

	min, max float64
}

// Description describes the validation in plain text formatting.
func (validator betweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be between %f and %f", validator.min, validator.max)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator betweenValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator betweenValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	f, ok := validateFloat(ctx, request, response)
	if !ok {
		return
	}

	if f < validator.min || f > validator.max {
		response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			request.AttributePath,
			"Invalid value",
			fmt.Sprintf("expected value to be in the range [%f, %f], got %f", validator.min, validator.max, f),
		))

		return
	}
}

// Between returns a new float value between validator.
func Between(min, max float64) tfsdk.AttributeValidator {
	if min > max {
		return nil
	}

	return betweenValidator{
		min: min,
		max: max,
	}
}
