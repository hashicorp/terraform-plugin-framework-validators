package listvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.AttributeValidator = valuesAreValidator{}

// valuesAreValidator validates that each list member validates against each of the value validators.
type valuesAreValidator struct {
	valueValidators []tfsdk.AttributeValidator
}

// Description describes the validation in plain text formatting.
func (v valuesAreValidator) Description(ctx context.Context) string {
	var descriptions []string
	for _, validator := range v.valueValidators {
		descriptions = append(descriptions, validator.Description(ctx))
	}

	return fmt.Sprintf("value must satisfy all validations: %s", strings.Join(descriptions, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v valuesAreValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v valuesAreValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	elems, ok := validateList(ctx, req, resp)
	if !ok {
		return
	}

	for k, elem := range elems {
		request := tfsdk.ValidateAttributeRequest{
			AttributePath:   req.AttributePath.WithElementKeyInt(k),
			AttributeConfig: elem,
			Config:          req.Config,
		}

		for _, validator := range v.valueValidators {
			validator.Validate(ctx, request, resp)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}
}

func ValuesAre(valueValidators ...tfsdk.AttributeValidator) tfsdk.AttributeValidator {
	return valuesAreValidator{
		valueValidators: valueValidators,
	}
}