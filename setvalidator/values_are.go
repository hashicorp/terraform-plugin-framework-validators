package setvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/hashicorp/terraform-plugin-framework-validators/validatordiag"
)

var _ tfsdk.AttributeValidator = valuesAreValidator{}

// valuesAreValidator validates that each set member validates against each of the value validators.
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
	elems, ok := validateSet(ctx, req, resp)
	if !ok {
		return
	}

	for k, elem := range elems {
		value, err := elem.ToTerraformValue(ctx)
		if err != nil {
			resp.Diagnostics.Append(validatordiag.AttributeValueTerraformValueDiagnostic(
				req.AttributePath,
				fmt.Sprintf("element at index: %d cannot be converted to Terraform value", k),
				err.Error(),
			))
			return
		}

		request := tfsdk.ValidateAttributeRequest{
			AttributePath:   req.AttributePath.WithElementKeyValue(value),
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
