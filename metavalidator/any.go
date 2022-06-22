package metavalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = anyValidator{}

// anyValidator validates that value validates against at least one of the value validators.
type anyValidator struct {
	valueValidators []tfsdk.AttributeValidator
}

// Description describes the validation in plain text formatting.
func (v anyValidator) Description(ctx context.Context) string {
	var descriptions []string
	for _, validator := range v.valueValidators {
		descriptions = append(descriptions, validator.Description(ctx))
	}

	return fmt.Sprintf("value must satisfy at least one of the validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v anyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
// If the number of iterations (i.e., k + 1) is greater than the number of diagnostics in the response then
// at least one of the validations has passed and, we return without any diagnostics.
func (v anyValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	for k, validator := range v.valueValidators {
		validator.Validate(ctx, req, resp)
		if k+1 > len(resp.Diagnostics) {
			resp.Diagnostics = []diag.Diagnostic{}
			return
		}
	}
}

// Any returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Validates against at least one of the value validators.
func Any(valueValidators ...tfsdk.AttributeValidator) tfsdk.AttributeValidator {
	return anyValidator{
		valueValidators: valueValidators,
	}
}
