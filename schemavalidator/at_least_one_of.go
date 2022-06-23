package schemavalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/attributepath"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// atLeastOneOfAttributeValidator is the underlying struct implementing AtLeastOneOf.
type atLeastOneOfAttributeValidator struct {
	attrPaths []*tftypes.AttributePath
}

// AtLeastOneOf checks that of a set of *tftypes.AttributePath,
// including the attribute it's applied to, at least one attribute out of all specified is configured.
//
// The provided tftypes.AttributePath must be "absolute",
// and starting with top level attribute names.
func AtLeastOneOf(attributePaths ...*tftypes.AttributePath) tfsdk.AttributeValidator {
	return &atLeastOneOfAttributeValidator{attributePaths}
}

var _ tfsdk.AttributeValidator = (*atLeastOneOfAttributeValidator)(nil)

func (av atLeastOneOfAttributeValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av atLeastOneOfAttributeValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensure that at least one attribute from this collection is set: %q", av.attrPaths)
}

func (av atLeastOneOfAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	// Assemble a slice of paths, ensuring we don't repeat the attribute this validator is applied to
	var paths []*tftypes.AttributePath
	if attributepath.Contains(req.AttributePath, av.attrPaths...) {
		paths = av.attrPaths
	} else {
		paths = append(av.attrPaths, req.AttributePath)
	}

	for _, path := range paths {
		var v attr.Value
		diags := req.Config.GetAttribute(ctx, path, &v)
		res.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		if !v.IsNull() {
			return
		}
	}

	res.Diagnostics.Append(validatordiag.InvalidAttributeSchemaDiagnostic(
		req.AttributePath,
		fmt.Sprintf("At least one attribute out of %q must be specified", attributepath.JoinToString(paths...)),
	))
}
