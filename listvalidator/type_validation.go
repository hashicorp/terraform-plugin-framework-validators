package listvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateList ensures that the request contains a List value.
func validateList(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) ([]attr.Value, bool) {
	var l types.List

	diags := tfsdk.ValueAs(ctx, request.AttributeConfig, &l)

	if diags.HasError() {
		response.Diagnostics.Append(diags...)

		return nil, false
	}

	if l.Unknown || l.Null {
		return nil, false
	}

	return l.Elems, true
}
