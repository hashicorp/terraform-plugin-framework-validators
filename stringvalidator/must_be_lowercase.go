// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.String = mustBeLowercaseValidator{}
var _ function.StringParameterValidator = mustBeLowercaseValidator{}

type mustBeLowercaseValidator struct{}

func (v mustBeLowercaseValidator) Description(_ context.Context) string {
	return "value must be lowercase"
}

func (v mustBeLowercaseValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v mustBeLowercaseValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if value != strings.ToLower(value) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

func (v mustBeLowercaseValidator) ValidateParameterString(ctx context.Context, request function.StringParameterValidatorRequest, response *function.StringParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueString()

	if value != strings.ToLower(value) {
		response.Error = validatorfuncerr.InvalidParameterValueMatchFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			value,
		)
	}
}

// MustBeLowercase returns a validator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a string.
//   - Is lowercase, meaning it is unchanged by strings.ToLower. Characters
//     without case, such as digits, underscores and spaces, are allowed, and
//     the empty string is valid.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func MustBeLowercase() mustBeLowercaseValidator {
	return mustBeLowercaseValidator{}
}
