package f64

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateFloat ensures that the request contains a Float64 value.
func validateFloat(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) (float64, bool) {
	var n types.Float64

	diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &n)

	if diags.HasError() {
		response.Diagnostics = append(response.Diagnostics, diags...)

		return 0, false
	}

	if n.Unknown || n.Null {
		return 0, false
	}

	return n.Value, true
}
