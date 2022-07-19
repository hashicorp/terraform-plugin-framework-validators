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

var _ tfsdk.AttributeValidator = atMostSumOfValidator{}

// atMostSumOfValidator validates that an integer Attribute's value is at most the sum of one
// or more integer Attributes retrieved via the given path expressions.
type atMostSumOfValidator struct {
	attributesToSumPathExpressions path.Expressions
}

// Description describes the validation in plain text formatting.
func (av atMostSumOfValidator) Description(_ context.Context) string {
	var attributePaths []string
	for _, p := range av.attributesToSumPathExpressions {
		attributePaths = append(attributePaths, p.String())
	}

	return fmt.Sprintf("value must be at most sum of %s", strings.Join(attributePaths, " + "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (av atMostSumOfValidator) MarkdownDescription(ctx context.Context) string {
	return av.Description(ctx)
}

// Validate performs the validation.
func (av atMostSumOfValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	i, ok := validateInt(ctx, request, response)
	if !ok {
		return
	}

	// Ensure input path expressions resolution against the current attribute
	expressions := request.AttributePathExpression.MergeExpressions(av.attributesToSumPathExpressions...)

	// Sum the value of all the attributes involved, but only if they are all known.
	var sumOfAttribs int64
	for _, expression := range expressions {
		matchedPaths, diags := request.Config.PathMatches(ctx, expression)
		response.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		for _, mp := range matchedPaths {
			// If the user specifies the same attribute this validator is applied to,
			// also as part of the input, skip it
			if mp.Equal(request.AttributePath) {
				continue
			}

			var attribToSum types.Int64
			diags := request.Config.GetAttribute(ctx, mp, &attribToSum)
			response.Diagnostics.Append(diags...)

			// Collect all errors
			if diags.HasError() {
				continue
			}

			// Delay validation until all involved attribute have a known value
			if attribToSum.IsUnknown() {
				return
			}

			// Attribute is null, so it doesn't contribute to the sum
			if attribToSum.IsNull() {
				continue
			}

			sumOfAttribs += attribToSum.Value
		}
	}

	if i > sumOfAttribs {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.AttributePath,
			av.Description(ctx),
			fmt.Sprintf("%d", i),
		))

		return
	}
}

// AtMostSumOf returns an AttributeValidator which ensures that any configured
// attribute value:
//
//     - Is a number, which can be represented by a 64-bit integer.
//     - Is at most the sum of the given attributes retrieved via the given path expression(s).
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtMostSumOf(attributesToSumPathExpressions ...path.Expression) tfsdk.AttributeValidator {
	return atMostSumOfValidator{attributesToSumPathExpressions}
}
