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
		response.Diagnostics.Append(diags...)

		return nil, false
	}

	if s.IsUnknown() || s.IsNull() {
		return nil, false
	}

	return s.Elements(), true
}
