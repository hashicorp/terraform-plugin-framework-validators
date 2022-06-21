package listvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/hashicorp/terraform-plugin-framework-validators/validatordiag"
)

var _ tfsdk.AttributeValidator = sizeAtLeastValidator{}

// sizeAtLeastValidator validates that list contains at least min elements.
type sizeAtLeastValidator struct {
	min int
}

// Description describes the validation in plain text formatting.
func (v sizeAtLeastValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("list must contain at least %d elements", v.min)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v sizeAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v sizeAtLeastValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	elems, ok := validateList(ctx, req, resp)
	if !ok {
		return
	}

	if len(elems) < v.min {
		resp.Diagnostics.Append(validatordiag.AttributeValueDiagnostic(
			req.AttributePath,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))

		return
	}
}

func SizeAtLeast(min int) tfsdk.AttributeValidator {
	return sizeAtLeastValidator{
		min: min,
	}
}
