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

// requiredWithAttributeValidator is the underlying struct implementing RequiredWith.
type requiredWithAttributeValidator struct {
	attrPaths []*tftypes.AttributePath
}

// RequiredWith checks that a set of *tftypes.AttributePath,
// including the attribute it's applied to, are set simultaneously.
// This implements the validation logic declaratively within the tfsdk.Schema.
//
// The provided tftypes.AttributePath must be "absolute",
// and starting with top level attribute names.
func RequiredWith(attributePaths ...*tftypes.AttributePath) tfsdk.AttributeValidator {
	return &requiredWithAttributeValidator{attributePaths}
}

var _ tfsdk.AttributeValidator = (*requiredWithAttributeValidator)(nil)

func (av requiredWithAttributeValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av requiredWithAttributeValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that if an attribute is set, also these are set: %q", av.attrPaths)
}

func (av requiredWithAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	var v attr.Value
	res.Diagnostics.Append(tfsdk.ValueAs(ctx, req.AttributeConfig, &v)...)
	if res.Diagnostics.HasError() {
		return
	}

	for _, path := range av.attrPaths {
		// If the user specifies the same attribute this validator is applied to,
		// also as part of the input, skip it.
		if req.AttributePath.Equal(path) {
			continue
		}

		var o attr.Value
		diags := req.Config.GetAttribute(ctx, path, &o)
		res.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		if !v.IsNull() && o.IsNull() {
			res.Diagnostics.Append(validatordiag.InvalidAttributeSchemaDiagnostic(
				req.AttributePath,
				fmt.Sprintf("Attribute %q must be specified when %q is specified", attributepath.ToString(path), attributepath.ToString(req.AttributePath)),
			))
		}
	}
}
