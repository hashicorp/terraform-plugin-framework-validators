package int64

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateInt ensures that the request contains an Int64 value.
func validateInt(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) (int64, bool) {
	var n types.Int64

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
