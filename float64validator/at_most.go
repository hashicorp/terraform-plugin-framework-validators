package float64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	f, diags := validateFloat(ctx, validator, request)

	if diags.HasError() {
		response.Diagnostics.Append(diags...)

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

func validateFloat(ctx context.Context, validator tfsdk.AttributeValidator, request tfsdk.ValidateAttributeRequest) (float64, diag.Diagnostics) {
	var n types.Float64

	diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &n)

	if diags.HasError() {
		var n types.Number

		diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &n)

		if diags.HasError() {
			return 0, diags
		}

		if n.Unknown {
			return 0, []diag.Diagnostic{
				validatordiag.AttributeValueDiagnostic(
					request.AttributePath,
					validator.Description(ctx),
					"Unknown",
				),
			}
		}

		if n.Null {
			return 0, []diag.Diagnostic{
				validatordiag.AttributeValueDiagnostic(
					request.AttributePath,
					validator.Description(ctx),
					"Null",
				),
			}
		}

		f, _ := n.Value.Float64()

		return f, nil

	}

	if n.Unknown {
		return 0, []diag.Diagnostic{
			validatordiag.AttributeValueDiagnostic(
				request.AttributePath,
				validator.Description(ctx),
				"Unknown",
			),
		}
	}

	if n.Null {
		return 0, []diag.Diagnostic{
			validatordiag.AttributeValueDiagnostic(
				request.AttributePath,
				validator.Description(ctx),
				"Null",
			),
		}
	}

	return n.Value, nil
}
