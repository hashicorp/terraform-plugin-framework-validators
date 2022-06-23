package setvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

	for _, elem := range elems {
		value, err := elem.ToTerraformValue(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Attribute Conversion Error During Set Element Validation",
				"An unexpected error was encountered when handling the a Set element. "+
					"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
					"Please report this to the provider developer:\n\n"+
					"Attribute Conversion Error During Set Element Validation.",
			)
			return
		}

		request := tfsdk.ValidateAttributeRequest{
			AttributePath:   req.AttributePath.WithElementKeyValue(value),
			AttributeConfig: elem,
			Config:          req.Config,
		}

		for _, validator := range v.valueValidators {
			validator.Validate(ctx, request, resp)
		}
	}
}

// ValuesAre returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Is a Set.
//     - Contains Set elements, each of which validate against each value validator.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func ValuesAre(valueValidators ...tfsdk.AttributeValidator) tfsdk.AttributeValidator {
	return valuesAreValidator{
		valueValidators: valueValidators,
	}
}