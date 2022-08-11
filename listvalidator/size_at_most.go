package listvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = sizeAtMostValidator{}

// sizeAtMostValidator validates that list contains at most max elements.
type sizeAtMostValidator struct {
	max int
}

// Description describes the validation in plain text formatting.
func (v sizeAtMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("list must contain at most %d elements", v.max)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v sizeAtMostValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v sizeAtMostValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	elems, ok := validateList(ctx, req, resp)
	if !ok {
		return
	}

	if len(elems) > v.max {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.AttributePath,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))

		return
	}
}

// SizeAtMost returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a List.
//   - Contains at most max elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeAtMost(max int) tfsdk.AttributeValidator {
	return sizeAtMostValidator{
		max: max,
	}
}
