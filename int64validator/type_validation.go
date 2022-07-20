package int64validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateInt ensures that the request contains an Int64 value.
func validateInt(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) (int64, bool) {
	var i types.Int64
	diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &i)
	response.Diagnostics.Append(diags...)
	if diags.HasError() {
		return 0, false
	}

	if i.IsUnknown() || i.IsNull() {
		return 0, false
	}

	return i.Value, true
}
