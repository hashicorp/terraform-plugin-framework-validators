// Copyright (c) HashiCorp, Inc.
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

var _ validator.Number = oneOfValidator{}
var _ function.NumberParameterValidator = oneOfValidator{}

type oneOfValidator struct {
	values []types.Number
}

func (v oneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v oneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be one of: %q", v.values)
}

func (v oneOfValidator) ValidateNumber(ctx context.Context, request validator.NumberRequest, response *validator.NumberResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, otherValue := range v.values {
		if value.Equal(otherValue) {
			return
		}
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
		request.Path,
		v.Description(ctx),
		value.String(),
	))
}

func (v oneOfValidator) ValidateParameterNumber(ctx context.Context, request function.NumberParameterValidatorRequest, response *function.NumberParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value

	for _, otherValue := range v.values {
		if value.Equal(otherValue) {
			return
		}
	}

	response.Error = validatorfuncerr.InvalidParameterValueMatchFuncError(
		request.ArgumentPosition,
		v.Description(ctx),
		value.String(),
	)
}

// OneOf checks that the Number held in the attribute or function parameter
// is one of the given `values`.
func OneOf(values ...*big.Float) oneOfValidator {
	frameworkValues := make([]types.Number, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.NumberValue(value))
	}

	return oneOfValidator{
		values: frameworkValues,
	}
}
