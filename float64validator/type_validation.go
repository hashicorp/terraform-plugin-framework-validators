package float64validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateFloat ensures that the request contains a Float64 value.
func validateFloat(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) (float64, bool) {
	t := request.AttributeConfig.Type(ctx)
	if t != types.Float64Type {
		response.Diagnostics.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			request.AttributePath,
			"expected value of type float64",
			t.String(),
		))
		return 0.0, false
	}

	f := request.AttributeConfig.(types.Float64)

	if f.IsUnknown() || f.IsNull() {
		return 0.0, false
	}

	return f.ValueFloat64(), true
}
