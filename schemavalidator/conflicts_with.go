package schemavalidator

import (
	"context"
	"fmt"

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
// including the attribute it's applied to, do not have a value simultaneously.
//
// This implements the validation logic declaratively within the tfsdk.Schema.
// Refer to [datasourcevalidator.Conflicting],
// [providervalidator.Conflicting], or [resourcevalidator.Conflicting]
// for declaring this type of validation outside the schema definition.
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
	// If attribute configuration is null, it cannot conflict with others
	if req.AttributeConfig.IsNull() {
		return
	}

	expressions := req.AttributePathExpression.MergeExpressions(av.pathExpressions...)

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
				res.Diagnostics.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
					req.AttributePath,
					fmt.Sprintf("Attribute %q cannot be specified when %q is specified", mp, req.AttributePath),
				))
			}
		}
	}
}
