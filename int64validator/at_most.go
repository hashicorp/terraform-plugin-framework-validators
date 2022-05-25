package int64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = atMostValidator{}

// atMostValidator validates that an integer Attribute's value is at most a certain value.
type atMostValidator struct {
	max int64
}

// Description describes the validation in plain text formatting.
func (validator atMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at most %d", validator.max)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator atMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator atMostValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	i, ok := validateInt(ctx, request, response)

	if !ok {
		return
	}

	if i > validator.max {
		response.Diagnostics.Append(validatordiag.AttributeValueDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			fmt.Sprintf("%d", i),
		))

		return
	}
}

// AtMost returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Is a number, which can be represented by a 64-bit integer.
//     - Is exclusively less than the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtMost(max int64) tfsdk.AttributeValidator {
	return atMostValidator{
		max: max,
	}
}
