package int64validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateInt ensures that the request contains an Int64 value.
func validateInt(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) (int64, bool) {
	t := request.AttributeConfig.Type(ctx)
	if t != types.Int64Type {
		response.Diagnostics.Append(validatordiag.InvalidTypeDiagnostic(
			request.AttributePath,
			"Expected value of type int64",
			t.String(),
		))
		return 0, false
	}

	i := request.AttributeConfig.(types.Int64)

	if i.Unknown || i.Null {
		return 0, false
	}

	return i.Value, true
}
