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

var _ validator.String = mustBeUppercaseValidator{}
var _ function.StringParameterValidator = mustBeUppercaseValidator{}

type mustBeUppercaseValidator struct{}

func (v mustBeUppercaseValidator) Description(_ context.Context) string {
	return "value must be uppercase"
}

func (v mustBeUppercaseValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v mustBeUppercaseValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if value != strings.ToUpper(value) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

func (v mustBeUppercaseValidator) ValidateParameterString(ctx context.Context, request function.StringParameterValidatorRequest, response *function.StringParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueString()

	if value != strings.ToUpper(value) {
		response.Error = validatorfuncerr.InvalidParameterValueMatchFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			value,
		)
	}
}

// MustBeUppercase returns a validator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a string.
//   - Is uppercase, meaning it is unchanged by strings.ToUpper. Characters
//     without case, such as digits, underscores and spaces, are allowed, and
//     the empty string is valid.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func MustBeUppercase() mustBeUppercaseValidator {
	return mustBeUppercaseValidator{}
}
