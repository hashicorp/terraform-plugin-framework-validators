package float64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		response.Diagnostics.Append(validatordiag.AttributeValueDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			fmt.Sprintf("%f", f),
		))

		return
	}
}

// AtMost returns a new float value at most validator.
func AtMost(max float64) tfsdk.AttributeValidator {
	return atMostValidator{
		max: max,
	}
}

func validateFloat(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) (float64, bool) {
	var n types.Float64

	diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &n)

	if diags.HasError() {
		var n types.Number

		diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &n)

		if diags.HasError() {
			response.Diagnostics = append(response.Diagnostics, diags...)

			return 0, false
		} else {
			if n.Unknown || n.Null {
				return 0, false
			}

			f, _ := n.Value.Float64()

			return f, true
		}
	} else {
		if n.Unknown || n.Null {
			return 0, false
		}

		return n.Value, true
	}
}
