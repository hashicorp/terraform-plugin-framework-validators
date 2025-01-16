// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int32validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.Int32 = atMostValidator{}
var _ function.Int32ParameterValidator = atMostValidator{}

type atMostValidator struct {
	max int32
}

func (validator atMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at most %d", validator.max)
}

func (validator atMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v atMostValidator) ValidateInt32(ctx context.Context, request validator.Int32Request, response *validator.Int32Response) {
	if request.ConfigValue.IsNull() {
		return
	}

	if request.ConfigValue.IsUnknown() {
		if refn, ok := request.ConfigValue.LowerBoundRefinement(); ok {
			if refn.IsInclusive() && refn.LowerBound() > v.max {
				response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
					request.Path,
					v.Description(ctx),
					fmt.Sprintf("unknown value that will be at least %d", refn.LowerBound()),
				))
			} else if !refn.IsInclusive() && refn.LowerBound() >= v.max {
				response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
					request.Path,
					v.Description(ctx),
					fmt.Sprintf("unknown value that will be greater than %d", refn.LowerBound()),
				))
			}
		}
		return
	}

	if request.ConfigValue.ValueInt32() > v.max {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", request.ConfigValue.ValueInt32()),
		))
	}
}

func (v atMostValidator) ValidateParameterInt32(ctx context.Context, request function.Int32ParameterValidatorRequest, response *function.Int32ParameterValidatorResponse) {
	if request.Value.IsNull() {
		return
	}

	if request.Value.IsUnknown() {
		if refn, ok := request.Value.LowerBoundRefinement(); ok {
			if refn.IsInclusive() && refn.LowerBound() > v.max {
				response.Error = validatorfuncerr.InvalidParameterValueFuncError(
					request.ArgumentPosition,
					v.Description(ctx),
					fmt.Sprintf("unknown value that will be at least %d", refn.LowerBound()),
				)
			} else if !refn.IsInclusive() && refn.LowerBound() >= v.max {
				response.Error = validatorfuncerr.InvalidParameterValueFuncError(
					request.ArgumentPosition,
					v.Description(ctx),
					fmt.Sprintf("unknown value that will be greater than %d", refn.LowerBound()),
				)
			}
		}
		return
	}

	if request.Value.ValueInt32() > v.max {
		response.Error = validatorfuncerr.InvalidParameterValueFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", request.Value.ValueInt32()),
		)
	}
}

// AtMost returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a number, which can be represented by a 32-bit integer.
//   - Is less than or equal to the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtMost(maxVal int32) atMostValidator {
	return atMostValidator{
		max: maxVal,
	}
}
