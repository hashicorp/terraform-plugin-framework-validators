package str

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateString ensures that the request contains a String value.
func validateString(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) (string, bool) {
	var s types.String

	diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &s)

	if diags.HasError() {
		response.Diagnostics = append(response.Diagnostics, diags...)

		return "", false
	}

	if s.Unknown || s.Null {
		return "", false
	}

	return s.Value, true
}
