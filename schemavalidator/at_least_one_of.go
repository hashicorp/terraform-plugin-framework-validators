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

// atLeastOneOfAttributeValidator is the underlying struct implementing AtLeastOneOf.
type atLeastOneOfAttributeValidator struct {
	pathExpressions path.Expressions
}

// AtLeastOneOf checks that of a set of path.Expression,
// including the attribute it's applied to,
// at least one attribute out of all specified is has a non-null value.
//
// This implements the validation logic declaratively within the tfsdk.Schema.
//
// Any relative path.Expression will be resolved against the attribute with this validator.
func AtLeastOneOf(attributePaths ...path.Expression) tfsdk.AttributeValidator {
	return &atLeastOneOfAttributeValidator{attributePaths}
}

var _ tfsdk.AttributeValidator = (*atLeastOneOfAttributeValidator)(nil)

func (av atLeastOneOfAttributeValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av atLeastOneOfAttributeValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that at least one attribute from this collection is set: %s", av.pathExpressions)
}

func (av atLeastOneOfAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	// If attribute configuration is not null, validator already succeeded.
	if !req.AttributeConfig.IsNull() {
		return
	}

	matchingPaths, diags := pathutils.PathMatchExpressionsAgainstAttributeConfig(ctx, av.pathExpressions, req.AttributePathExpression, req.Config)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Validate values at the matching paths
	for _, mp := range matchingPaths {
		var mpVal attr.Value
		diags := req.Config.GetAttribute(ctx, mp, &mpVal)
		res.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		// Delay validation until all involved attribute have a known value
		if mpVal.IsUnknown() {
			return
		}

		if !mpVal.IsNull() {
			return
		}
	}

	res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
		req.AttributePath,
		fmt.Sprintf("At least one attribute out of %q must be specified", matchingPaths),
	))
}
