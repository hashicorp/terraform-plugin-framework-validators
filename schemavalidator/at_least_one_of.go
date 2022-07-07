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
// including the attribute it's applied to, at least one attribute out of all specified is configured.
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
	matchingPaths, diags := pathutils.PathMatchExpressionsAgainstAttributeConfig(ctx, av.pathExpressions, req.AttributePathExpression, req.Config)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Validate values at the matching paths
	for _, p := range matchingPaths {
		var v attr.Value
		diags := req.Config.GetAttribute(ctx, p, &v)
		res.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		if !v.IsNull() {
			return
		}
	}

	res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
		req.AttributePath,
		fmt.Sprintf("At least one attribute out of %q must be specified", matchingPaths),
	))
}
