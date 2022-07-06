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

// conflictsWithAttributeValidator is the underlying struct implementing ConflictsWith.
type conflictsWithAttributeValidator struct {
	pathExpressions path.Expressions
}

// ConflictsWith checks that a set of path.Expression,
// including the attribute it's applied to, are not set simultaneously.
// This implements the validation logic declaratively within the tfsdk.Schema.
//
// Relative path.Expression will be resolved against the validated attribute.
func ConflictsWith(attributePaths ...path.Expression) tfsdk.AttributeValidator {
	return &conflictsWithAttributeValidator{attributePaths}
}

var _ tfsdk.AttributeValidator = (*conflictsWithAttributeValidator)(nil)

func (av conflictsWithAttributeValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av conflictsWithAttributeValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that if an attribute is set, these are not set: %q", av.pathExpressions)
}

func (av conflictsWithAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	var v attr.Value
	res.Diagnostics.Append(tfsdk.ValueAs(ctx, req.AttributeConfig, &v)...)
	if res.Diagnostics.HasError() {
		return
	}

	matchingPaths, diags := pathutils.PathMatchExpressionsAgainstAttributeConfig(ctx, av.pathExpressions, req.AttributePathExpression, req.Config)
	res.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Validate values at the matching paths
	for _, p := range matchingPaths {
		// If the user specifies the same attribute this validator is applied to,
		// also as part of the input, skip it.
		if p.Equal(req.AttributePath) {
			continue
		}

		var o attr.Value
		diags := req.Config.GetAttribute(ctx, p, &o)
		res.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		if !v.IsNull() && !o.IsNull() {
			res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
				req.AttributePath,
				fmt.Sprintf("Attribute %q cannot be specified when %q is specified", p, req.AttributePath),
			))
		}
	}
}
