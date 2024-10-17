// Copyright (c) HashiCorp, Inc.
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

var _ validator.Float32 = atMostValidator{}
var _ function.Float32ParameterValidator = atMostValidator{}

type atMostValidator struct {
	max float32
}

func (validator atMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at most %f", validator.max)
}

func (validator atMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v atMostValidator) ValidateFloat32(ctx context.Context, request validator.Float32Request, response *validator.Float32Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueFloat32()

	if value > v.max {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%f", value),
		))
	}
}

func (v atMostValidator) ValidateParameterFloat32(ctx context.Context, request function.Float32ParameterValidatorRequest, response *function.Float32ParameterValidatorResponse) {
	if request.Value.IsNull() || request.Value.IsUnknown() {
		return
	}

	value := request.Value.ValueFloat32()

	if value > v.max {
		response.Error = validatorfuncerr.InvalidParameterValueFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%f", value),
		)
	}
}

// AtMost returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a number, which can be represented by a 32-bit floating point.
//   - Is less than or equal to the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtMost(maxVal float32) atMostValidator {
	return atMostValidator{
		max: maxVal,
	}
}
