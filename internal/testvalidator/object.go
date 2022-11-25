package testvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Object = ObjectValidator{}

type ObjectValidator struct {
	Diagnostics diag.Diagnostics
}

func (v ObjectValidator) Description(ctx context.Context) string {
	return "returns given Diagnostics"
}

func (v ObjectValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ObjectValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	resp.Diagnostics = v.Diagnostics
}
