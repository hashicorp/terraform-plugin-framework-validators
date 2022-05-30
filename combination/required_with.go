package combination

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/attr_path"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/diag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// requiredWithValidator checks that a set of *tftypes.AttributePath,
// including the attribute it's applied to, are set simultaneously.
// This implements the validation logic declaratively within the tfsdk.Schema.
//
// The provided tftypes.AttributePath must be "absolute",
// and starting with top level attribute names.
type requiredWithValidator struct {
	attrPaths []*tftypes.AttributePath
}

// RequiredWith is a helper to instantiate requiredWithValidator.
func RequiredWith(attributePaths ...*tftypes.AttributePath) tfsdk.AttributeValidator {
	return &requiredWithValidator{attributePaths}
}

var _ tfsdk.AttributeValidator = (*requiredWithValidator)(nil)

func (av requiredWithValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av requiredWithValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that if an attribute is set, also these are set: %q", av.attrPaths)
}

func (av requiredWithValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	var v attr.Value
	res.Diagnostics.Append(tfsdk.ValueAs(ctx, req.AttributeConfig, &v)...)
	if res.Diagnostics.HasError() {
		return
	}

	for _, path := range av.attrPaths {
		var o attr.Value
		res.Diagnostics.Append(req.Config.GetAttribute(ctx, path, &o)...)
		if res.Diagnostics.HasError() {
			return
		}

		if !v.IsNull() && o.IsNull() {
			res.Diagnostics.Append(diag.InvalidCombinationDiagnostic(
				req.AttributePath,
				fmt.Sprintf("Attribute %q must be specified when %q is specified", attr_path.ToString(path), attr_path.ToString(req.AttributePath)),
			))
			return
		}
	}
}
