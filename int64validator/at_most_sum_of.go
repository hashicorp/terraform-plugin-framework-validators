package int64validator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/validatordiag"
)

var _ tfsdk.AttributeValidator = atMostSumOfValidator{}

// atMostSumOfValidator validates that an integer Attribute's value is at most the sum of one
// or more integer Attributes.
type atMostSumOfValidator struct {
	attributesToSumPaths []*tftypes.AttributePath
}

// Description describes the validation in plain text formatting.
func (validator atMostSumOfValidator) Description(_ context.Context) string {
	var attributePaths []string
	for _, path := range validator.attributesToSumPaths {
		attributePaths = append(attributePaths, path.String())
	}

	return fmt.Sprintf("value must be at most sum of %s", strings.Join(attributePaths, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator atMostSumOfValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
func (validator atMostSumOfValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	i, ok := validateInt(ctx, request, response)

	if !ok {
		return
	}

	var sumOfAttribs int64

	for _, path := range validator.attributesToSumPaths {
		var attribToSum types.Int64

		response.Diagnostics.Append(request.Config.GetAttribute(ctx, path, &attribToSum)...)
		if response.Diagnostics.HasError() {
			return
		}

		sumOfAttribs += attribToSum.Value
	}

	if i > sumOfAttribs {

		response.Diagnostics.Append(validatordiag.AttributeValueDiagnostic(
			request.AttributePath,
			validator.Description(ctx),
			fmt.Sprintf("%d", i),
		))

		return
	}
}

// AtMostSumOf returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Is a number, which can be represented by a 64-bit integer.
//     - Is exclusively at most the sum of the given attributes.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtMostSumOf(attributesToSum []*tftypes.AttributePath) tfsdk.AttributeValidator {
	return atMostSumOfValidator{
		attributesToSumPaths: attributesToSum,
	}
}
