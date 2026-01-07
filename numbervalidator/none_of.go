// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package numbervalidator

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.Number = noneOfValidator{}
var _ function.NumberParameterValidator = noneOfValidator{}

type noneOfValidator struct {
	values []types.Number
}

func (v noneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v noneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be none of: %q", v.values)
}

func (v noneOfValidator) ValidateNumber(ctx context.Context, request validator.NumberRequest, response *validator.NumberResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, otherValue := range v.values {
		if !value.Equal(otherValue) {
			continue
		}

		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value.String(),
		))

		break
	}
}

func (v noneOfValidator) ValidateParameterNumber(ctx context.Context, request function.NumberParameterValidatorRequest, response *function.NumberParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value

	for _, otherValue := range v.values {
		if !value.Equal(otherValue) {
			continue
		}

		response.Error = validatorfuncerr.InvalidParameterValueMatchFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			value.String(),
		)

		break
	}
}

// NoneOf checks that the Number held in the attribute or function parameter
// is none of the given `values`.
func NoneOf(values ...*big.Float) noneOfValidator {
	frameworkValues := make([]types.Number, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.NumberValue(value))
	}

	return noneOfValidator{
		values: frameworkValues,
	}
}
