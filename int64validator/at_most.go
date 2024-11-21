// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.Int64 = atMostValidator{}
var _ function.Int64ParameterValidator = atMostValidator{}

type atMostValidator struct {
	max int64
}

func (validator atMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at most %d", validator.max)
}

func (validator atMostValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (v atMostValidator) ValidateInt64(ctx context.Context, request validator.Int64Request, response *validator.Int64Response) {
	if request.ConfigValue.IsNull() {
		return
	}

	if request.ConfigValue.IsUnknown() {
		// Check if there is a lower bound refinement, and if that lower bound indicates the eventual value will be invalid
		if lowerRefn, ok := request.ConfigValue.LowerBoundRefinement(); ok {
			if lowerRefn.IsInclusive() && lowerRefn.LowerBound() > v.max {
				response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
					request.Path,
					"Invalid Attribute Value",
					// TODO: improve error messaging?
					fmt.Sprintf("Attribute %s %s, got an unknown value that will be greater than: %d", request.Path, v.Description(ctx), lowerRefn.LowerBound()),
				))
			} else if !lowerRefn.IsInclusive() && lowerRefn.LowerBound() >= v.max {
				response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
					request.Path,
					"Invalid Attribute Value",
					fmt.Sprintf("Attribute %s %s, got an unknown value that will be at least: %d", request.Path, v.Description(ctx), lowerRefn.LowerBound()),
				))
			}
		}
		return
	}

	if request.ConfigValue.ValueInt64() > v.max {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", request.ConfigValue.ValueInt64()),
		))
	}
}

func (v atMostValidator) ValidateParameterInt64(ctx context.Context, request function.Int64ParameterValidatorRequest, response *function.Int64ParameterValidatorResponse) {
	if request.Value.IsNull() {
		return
	}

	if request.Value.IsUnknown() {
		// Check if there is a lower bound refinement, and if that lower bound indicates the eventual value will be invalid
		if lowerRefn, ok := request.Value.LowerBoundRefinement(); ok {
			if lowerRefn.IsInclusive() && lowerRefn.LowerBound() > v.max {
				response.Error = function.NewArgumentFuncError(
					request.ArgumentPosition,
					fmt.Sprintf("Invalid Parameter Value: %s, got an unknown value that will be greater than: %d", v.Description(ctx), lowerRefn.LowerBound()),
				)
			} else if !lowerRefn.IsInclusive() && lowerRefn.LowerBound() >= v.max {
				response.Error = function.NewArgumentFuncError(
					request.ArgumentPosition,
					fmt.Sprintf("Invalid Parameter Value: %s, got an unknown value that will be at least: %d", v.Description(ctx), lowerRefn.LowerBound()),
				)
			}
		}
		return
	}

	if request.Value.ValueInt64() > v.max {
		response.Error = validatorfuncerr.InvalidParameterValueFuncError(
			request.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", request.Value.ValueInt64()),
		)
	}
}

// AtMost returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a number, which can be represented by a 64-bit integer.
//   - Is less than or equal to the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtMost(maxVal int64) atMostValidator {
	return atMostValidator{
		max: maxVal,
	}
}
