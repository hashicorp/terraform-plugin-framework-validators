package float64

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = atMostValidator{}

// atMostValidator validates that an float Attribute's value is at most a certain value.
type atMostValidator struct {
	max float64
}

// Description describes the validation in plain text formatting.
func (validator atMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at most %f", validator.max)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator atMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator atMostValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	f, ok := validateFloat(ctx, request, response)

	if !ok {
		return
	}

	if f > validator.max {
		response.Diagnostics.Append(diag.InvalidValueDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			fmt.Sprintf("%f", f),
		))

		return
	}
}

// AtMost returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Is a number, which can be represented by a 64-bit floating point.
//     - Is exclusively less than the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtMost(max float64) tfsdk.AttributeValidator {
	return atMostValidator{
		max: max,
	}
}
