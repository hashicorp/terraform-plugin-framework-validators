// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package float32validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.Float32 = atLeastValidator{}
var _ function.Float32ParameterValidator = atLeastValidator{}

type atLeastValidator struct {
	min float32
}

func (validator atLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at least %f", validator.min)
}

func (validator atLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator atLeastValidator) ValidateFloat32(ctx context.Context, request validator.Float32Request, response *validator.Float32Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueFloat32()

	if value < validator.min {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			validator.Description(ctx),
			fmt.Sprintf("%f", value),
		))
	}
}

func (validator atLeastValidator) ValidateParameterFloat32(ctx context.Context, request function.Float32ParameterValidatorRequest, response *function.Float32ParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueFloat32()

	if value < validator.min {
		response.Error = validatorfuncerr.InvalidParameterValueFuncError(
			request.ArgumentPosition,
			validator.Description(ctx),
			fmt.Sprintf("%f", value),
		)
	}
}

// AtLeast returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a number, which can be represented by a 32-bit floating point.
//   - Is greater than or equal to the given minimum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtLeast(minVal float32) atLeastValidator {
	return atLeastValidator{
		min: minVal,
	}
}
