// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float32validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.Float32 = noneOfValidator{}
var _ function.Float32ParameterValidator = noneOfValidator{}

type noneOfValidator struct {
	values []types.Float32
}

func (v noneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v noneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be none of: %q", v.values)
}

func (v noneOfValidator) ValidateFloat32(ctx context.Context, request validator.Float32Request, response *validator.Float32Response) {
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

func (v noneOfValidator) ValidateParameterFloat32(ctx context.Context, request function.Float32ParameterValidatorRequest, response *function.Float32ParameterValidatorResponse) {
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

// NoneOf checks that the float32 held in the attribute or function parameter
// is none of the given `values`.
func NoneOf(values ...float32) noneOfValidator {
	frameworkValues := make([]types.Float32, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.Float32Value(value))
	}

	return noneOfValidator{
		values: frameworkValues,
	}
}
