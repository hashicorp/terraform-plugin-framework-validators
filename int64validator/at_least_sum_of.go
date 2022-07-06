package int64validator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ tfsdk.AttributeValidator = atLeastSumOfValidator{}

// atLeastSumOfValidator validates that an integer Attribute's value is at least the sum of one
// or more integer Attributes.
type atLeastSumOfValidator struct {
	attributesToSumPaths []path.Path
}

// Description describes the validation in plain text formatting.
func (validator atLeastSumOfValidator) Description(_ context.Context) string {
	var attributePaths []string
	for _, p := range validator.attributesToSumPaths {
		attributePaths = append(attributePaths, p.String())
	}

	return fmt.Sprintf("value must be at least sum of %s", strings.Join(attributePaths, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator atLeastSumOfValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator atLeastSumOfValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	i, ok := validateInt(ctx, request, response)

	if !ok {
		return
	}

	var sumOfAttribs int64
	var numUnknownAttribsToSum int

	for _, p := range validator.attributesToSumPaths {
		var attribToSum types.Int64

		response.Diagnostics.Append(request.Config.GetAttribute(ctx, p, &attribToSum)...)
		if response.Diagnostics.HasError() {
			return
		}

		if attribToSum.Null {
			continue
		}

		if attribToSum.Unknown {
			numUnknownAttribsToSum++
			continue
		}

		sumOfAttribs += attribToSum.Value
	}

	if numUnknownAttribsToSum == len(validator.attributesToSumPaths) {
		return
	}

	if i < sumOfAttribs {

		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			fmt.Sprintf("%d", i),
		))

		return
	}
}

// AtLeastSumOf returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Is a number, which can be represented by a 64-bit integer.
//     - Is exclusively at least the sum of the given attributes.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtLeastSumOf(attributesToSum ...path.Path) tfsdk.AttributeValidator {
	return atLeastSumOfValidator{
		attributesToSumPaths: attributesToSum,
	}
}
