package pathutils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// PathMatchExpressionsAgainstAttributeConfig returns the path.Paths matching the given path.Expressions.
//
// Each path.Expression has been merged with the given attribute path.Expression
// (likely from the tfsdk.ValidateAttributeRequest), resolved,
// and then matched against the given attribute tfsdk.Config (also from the tfsdk.ValidateAttributeRequest).
//
// This is useful for tfsdk.AttributeValidator that accept path.Expressions, and validate the attributes matching
// to the expressions, in relation to the attribute the validator is applied to.
// For example usage, please look at the `schemavalidator` package in this repository.
func PathMatchExpressionsAgainstAttributeConfig(ctx context.Context, pathExps path.Expressions, attrPathExp path.Expression, attrConfig tfsdk.Config) (path.Paths, diag.Diagnostics) {
	var resDiags diag.Diagnostics

	pathExpressions := MergeExpressionsWithAttribute(pathExps, attrPathExp)

	resPaths := make(path.Paths, 0, len(pathExpressions))

	for _, pe := range pathExpressions {
		// Retrieve all the attribute paths that match the given expressions
		matchingPaths, diags := attrConfig.PathMatches(ctx, pe)
		resDiags.Append(diags...)
		if diags.HasError() {
			return nil, resDiags
		}

		// Confirm at least one attribute was matched.
		// If not, collect errors so that the callee can bubble the bugs up.
		if len(matchingPaths) == 0 {
			resDiags.Append(validatordiag.BugInProviderDiagnostic(fmt.Sprintf("Path expression %q matches no attribute", pe)))
		}

		resPaths.Append(matchingPaths...)
	}

	return resPaths, resDiags
}
