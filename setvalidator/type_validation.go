package setvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateSet ensures that the request contains a Set value.
func validateSet(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) ([]attr.Value, bool) {
	var n types.Set

	diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &n)

	if diags.HasError() {
		response.Diagnostics = append(response.Diagnostics, diags...)

		return nil, false
	}

	if n.Unknown || n.Null {
		return nil, false
	}

	return n.Elems, true
}
