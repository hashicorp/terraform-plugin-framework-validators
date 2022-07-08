package schemavalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/pathutils"
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
// including the attribute it's applied to, one and only one attribute out of all specified is configured.
// It will also cause a validation error if none are specified.
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
	matchingPaths, diags := pathutils.PathMatchExpressionsAgainstAttributeConfig(ctx, av.pathExpressions, req.AttributePathExpression, req.Config)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Validate values at the matching paths
	count := 0
	for _, mp := range matchingPaths {
		var mpVal attr.Value
		diags := req.Config.GetAttribute(ctx, mp, &mpVal)
		res.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		// Delay validation until all involved attribute
		// have a known value
		if mpVal.IsUnknown() {
			return
		}

		if !mpVal.IsNull() {
			count++
		}
	}

	if count == 0 {
		res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			req.AttributePath,
			fmt.Sprintf("No attribute specified when one (and only one) of %q is required", matchingPaths),
		))
	}

	if count > 1 {
		res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			req.AttributePath,
			fmt.Sprintf("%d attributes specified when one (and only one) of %q is required", count, matchingPaths),
		))
	}
}
