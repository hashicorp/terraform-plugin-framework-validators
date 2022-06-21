package setvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateSet ensures that the request contains a Set value.
func validateSet(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) ([]attr.Value, bool) {
	var s types.Set

	diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &s)

	if diags.HasError() {
		response.Diagnostics = append(response.Diagnostics, diags...)

		return nil, false
	}

	if s.Unknown || s.Null {
		return nil, false
	}

	return s.Elems, true
}
