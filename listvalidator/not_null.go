package listvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.List = notNullValidator{}

// notNullValidator validates that list is not null
type notNullValidator struct{}

// Description describes the validation in plain text formatting.
func (v notNullValidator) Description(_ context.Context) string {
	return "list must not be null"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v notNullValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v notNullValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeDiagnostic(
			req.Path,
			v.Description(ctx),
		))
	}
}

// NotNull returns a validator which ensures that any configured list is set
// to a value.
//
// This validator is equivalent to the `Required` field on attributes and is only
// practical for use with `schema.ListNestedBlock`
func NotNull() validator.List {
	return notNullValidator{}
}
