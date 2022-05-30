package float64

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = betweenValidator{}

// betweenValidator validates that an float Attribute's value is in a range.
type betweenValidator struct {
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
		response.Diagnostics.Append(diag.InvalidValueDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			fmt.Sprintf("%f", f),
		))

		return
	}
}

// Between returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Is a number, which can be represented by a 64-bit floating point.
//     - Is exclusively greater than the given minimum and less than the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func Between(min, max float64) tfsdk.AttributeValidator {
	if min > max {
		return nil
	}

	return betweenValidator{
		min: min,
		max: max,
	}
}
