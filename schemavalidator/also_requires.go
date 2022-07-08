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

// alsoRequiresAttributeValidator is the underlying struct implementing AlsoRequires.
type alsoRequiresAttributeValidator struct {
	pathExpressions path.Expressions
}

// AlsoRequires checks that a set of path.Expression has a non-null value,
// if the current attribute also has a non-null value.
//
// This implements the validation logic declaratively within the tfsdk.Schema.
//
// Relative path.Expression will be resolved against the validated attribute.
func AlsoRequires(attributePaths ...path.Expression) tfsdk.AttributeValidator {
	return &alsoRequiresAttributeValidator{attributePaths}
}

var _ tfsdk.AttributeValidator = (*alsoRequiresAttributeValidator)(nil)

func (av alsoRequiresAttributeValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av alsoRequiresAttributeValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that if an attribute is set, also these are set: %q", av.pathExpressions)
}

func (av alsoRequiresAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	// If attribute configuration is null, there is nothing else to validate
	if req.AttributeConfig.IsNull() {
		return
	}

	matchingPaths, diags := pathutils.PathMatchExpressionsAgainstAttributeConfig(ctx, av.pathExpressions, req.AttributePathExpression, req.Config)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Validate values at the matching paths
	for _, mp := range matchingPaths {
		// If the user specifies the same attribute this validator is applied to,
		// also as part of the input, skip it.
		if mp.Equal(req.AttributePath) {
			continue
		}

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

		if mpVal.IsNull() {
			res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
				req.AttributePath,
				fmt.Sprintf("Attribute %q must be specified when %q is specified", mp, req.AttributePath),
			))
		}
	}
}
