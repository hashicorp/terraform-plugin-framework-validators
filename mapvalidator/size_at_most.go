package mapvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/hashicorp/terraform-plugin-framework-validators/validatordiag"
)

var _ tfsdk.AttributeValidator = sizeAtMostValidator{}

// sizeAtMostValidator validates that map contains at most max elements.
type sizeAtMostValidator struct {
	max int
}

// Description describes the validation in plain text formatting.
func (v sizeAtMostValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("map must contain at most %d elements", v.max)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v sizeAtMostValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v sizeAtMostValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	elems, ok := validateMap(ctx, req, resp)
	if !ok {
		return
	}

	if len(elems) > v.max {
		resp.Diagnostics.Append(validatordiag.AttributeValueDiagnostic(
			req.AttributePath,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))

		return
	}
}

func SizeAtMost(max int) tfsdk.AttributeValidator {
	return sizeAtMostValidator{
		max: max,
	}
}
