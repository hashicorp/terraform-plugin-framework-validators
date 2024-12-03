// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package setvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.Set = sizeBetweenValidator{}
var _ function.SetParameterValidator = sizeBetweenValidator{}

type sizeBetweenValidator struct {
	min int
	max int
}

func (v sizeBetweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("set must contain at least %d elements and at most %d elements", v.min, v.max)
}

func (v sizeBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v sizeBetweenValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() {
		return
	}

	if req.ConfigValue.IsUnknown() {
		if refn, ok := req.ConfigValue.LengthLowerBoundRefinement(); ok && refn.LowerBound() > int64(v.max) {
			resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
				req.Path,
				v.Description(ctx),
				fmt.Sprintf("unknown value that will have at least %d elements", refn.LowerBound()),
			))
		}

		if refn, ok := req.ConfigValue.LengthUpperBoundRefinement(); ok && refn.UpperBound() < int64(v.min) {
			resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
				req.Path,
				v.Description(ctx),
				fmt.Sprintf("unknown value that will have at most %d elements", refn.UpperBound()),
			))
		}
		return
	}

	elems := req.ConfigValue.Elements()

	if len(elems) < v.min || len(elems) > v.max {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))
	}
}

func (v sizeBetweenValidator) ValidateParameterSet(ctx context.Context, req function.SetParameterValidatorRequest, resp *function.SetParameterValidatorResponse) {
	if req.Value.IsNull() {
		return
	}

	if req.Value.IsUnknown() {
		if refn, ok := req.Value.LengthLowerBoundRefinement(); ok && refn.LowerBound() > int64(v.max) {
			resp.Error = validatorfuncerr.InvalidParameterValueFuncError(
				req.ArgumentPosition,
				v.Description(ctx),
				fmt.Sprintf("unknown value that will have at least %d elements", refn.LowerBound()),
			)
		}

		if refn, ok := req.Value.LengthUpperBoundRefinement(); ok && refn.UpperBound() < int64(v.min) {
			resp.Error = validatorfuncerr.InvalidParameterValueFuncError(
				req.ArgumentPosition,
				v.Description(ctx),
				fmt.Sprintf("unknown value that will have at most %d elements", refn.UpperBound()),
			)
		}
		return
	}

	elems := req.Value.Elements()

	if len(elems) < v.min || len(elems) > v.max {
		resp.Error = validatorfuncerr.InvalidParameterValueFuncError(
			req.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		)
	}
}

// SizeBetween returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a Set.
//   - Contains at least min elements and at most max elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeBetween(minVal, maxVal int) sizeBetweenValidator {
	return sizeBetweenValidator{
		min: minVal,
		max: maxVal,
	}
}
