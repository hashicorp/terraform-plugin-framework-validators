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

// exactlyOneOfAttributeValidator is the underlying struct implementing ExactlyOneOf.
type exactlyOneOfAttributeValidator struct {
	attrPaths []*tftypes.AttributePath
}

// ExactlyOneOf checks that of a set of *tftypes.AttributePath,
// including the attribute it's applied to, one and only one attribute out of all specified is configured.
// It will also cause a validation error if none are specified.
//
// The provided tftypes.AttributePath must be "absolute",
// and starting with top level attribute names.
func ExactlyOneOf(attributePaths ...*tftypes.AttributePath) tfsdk.AttributeValidator {
	return &exactlyOneOfAttributeValidator{attributePaths}
}

var _ tfsdk.AttributeValidator = (*exactlyOneOfAttributeValidator)(nil)

func (av exactlyOneOfAttributeValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av exactlyOneOfAttributeValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that one and only one attribute from this collection is set: %q", av.attrPaths)
}

func (av exactlyOneOfAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	// Assemble a slice of paths, ensuring we don't repeat the attribute this validator is applied to
	var paths []*tftypes.AttributePath
	if attributepath.Contains(req.AttributePath, av.attrPaths...) {
		paths = av.attrPaths
	} else {
		paths = append(av.attrPaths, req.AttributePath)
	}

	count := 0
	for _, path := range paths {
		var v attr.Value
		diags := req.Config.GetAttribute(ctx, path, &v)
		res.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		if !v.IsNull() {
			count++
		}
	}

	if count == 0 {
		res.Diagnostics.Append(validatordiag.InvalidAttributeSchemaDiagnostic(
			req.AttributePath,
			fmt.Sprintf("No attribute specified when one (and only one) of %q is required", attributepath.JoinToString(paths...)),
		))
	}

	if count > 1 {
		res.Diagnostics.Append(validatordiag.InvalidAttributeSchemaDiagnostic(
			req.AttributePath,
			fmt.Sprintf("%d attributes specified when one (and only one) of %q is required", count, attributepath.JoinToString(paths...)),
		))
	}
}
