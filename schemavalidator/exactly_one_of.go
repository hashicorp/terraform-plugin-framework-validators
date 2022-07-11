package schemavalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// exactlyOneOfAttributeValidator is the underlying struct implementing ExactlyOneOf.
type exactlyOneOfAttributeValidator struct {
	pathExpressions path.Expressions
}

// ExactlyOneOf checks that of a set of path.Expression,
// including the attribute it's applied to,
// one and only one attribute out of all specified has a value.
// It will also cause a validation error if none are specified.
//
// This implements the validation logic declaratively within the tfsdk.Schema.
//
// Relative path.Expression will be resolved against the validated attribute.
func ExactlyOneOf(attributePaths ...path.Expression) tfsdk.AttributeValidator {
	return &exactlyOneOfAttributeValidator{attributePaths}
}

var _ tfsdk.AttributeValidator = (*exactlyOneOfAttributeValidator)(nil)

func (av exactlyOneOfAttributeValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av exactlyOneOfAttributeValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that one and only one attribute from this collection is set: %q", av.pathExpressions)
}

func (av exactlyOneOfAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	count := 0
	expressions := req.AttributePathExpression.MergeExpressions(av.pathExpressions...)

	// If current attribute is unknown, delay validation
	if req.AttributeConfig.IsUnknown() {
		return
	}

	// Now that we know the current attribute is known, check whether it is
	// null to determine if it should contribute to the count. Later logic
	// will remove a duplicate matching path, should it be included in the
	// given expressions.
	if !req.AttributeConfig.IsNull() {
		count++
	}

	for _, expression := range expressions {
		matchedPaths, diags := req.Config.PathMatches(ctx, expression)

		res.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		for _, mp := range matchedPaths {
			// If the user specifies the same attribute this validator is applied to,
			// also as part of the input, skip it
			if mp.Equal(req.AttributePath) {
				continue
			}

			var mpVal attr.Value
			diags := req.Config.GetAttribute(ctx, mp, &mpVal)
			res.Diagnostics.Append(diags...)

			// Collect all errors
			if diags.HasError() {
				continue
			}

			// Delay validation until all involved attribute have a known value
			if mpVal.IsUnknown() {
				return
			}

			if !mpVal.IsNull() {
				count++
			}
		}
	}

	if count == 0 {
		res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			req.AttributePath,
			fmt.Sprintf("No attribute specified when one (and only one) of %s is required", expressions),
		))
	}

	if count > 1 {
		res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			req.AttributePath,
			fmt.Sprintf("%d attributes specified when one (and only one) of %s is required", count, expressions),
		))
	}
}
