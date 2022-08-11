package setvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = sizeBetweenValidator{}

// sizeBetweenValidator validates that set contains at least min elements
// and at most max elements.
type sizeBetweenValidator struct {
	min int
	max int
}

// Description describes the validation in plain text formatting.
func (v sizeBetweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("set must contain at least %d elements and at most %d elements", v.min, v.max)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v sizeBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v sizeBetweenValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	elems, ok := validateSet(ctx, req, resp)
	if !ok {
		return
	}

	if len(elems) < v.min || len(elems) > v.max {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.AttributePath,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))

		return
	}
}

// SizeBetween returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a Set.
//   - Contains at least min elements and at most max elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeBetween(min, max int) tfsdk.AttributeValidator {
	return sizeBetweenValidator{
		min: min,
		max: max,
	}
}
