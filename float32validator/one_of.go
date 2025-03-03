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

var _ validator.Float32 = oneOfValidator{}
var _ function.Float32ParameterValidator = oneOfValidator{}

type oneOfValidator struct {
	values []types.Float32
}

func (v oneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v oneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be one of: %q", v.values)
}

func (v oneOfValidator) ValidateFloat32(ctx context.Context, request validator.Float32Request, response *validator.Float32Response) {
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

func (v oneOfValidator) ValidateParameterFloat32(ctx context.Context, request function.Float32ParameterValidatorRequest, response *function.Float32ParameterValidatorResponse) {
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

// OneOf checks that the float32 held in the attribute or function parameter
// is one of the given `values`.
func OneOf(values ...float32) oneOfValidator {
	frameworkValues := make([]types.Float32, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.Float32Value(value))
	}

	return oneOfValidator{
		values: frameworkValues,
	}
}
