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

var _ validator.Set = sizeAtLeastValidator{}
var _ function.SetParameterValidator = sizeAtLeastValidator{}

type sizeAtLeastValidator struct {
	min int
}

func (v sizeAtLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("set must contain at least %d elements", v.min)
}

func (v sizeAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v sizeAtLeastValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() {
		return
	}

	if req.ConfigValue.IsUnknown() {
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

	if len(elems) < v.min {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))
	}
}

func (v sizeAtLeastValidator) ValidateParameterSet(ctx context.Context, req function.SetParameterValidatorRequest, resp *function.SetParameterValidatorResponse) {
	if req.Value.IsNull() {
		return
	}

	if req.Value.IsUnknown() {
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

	if len(elems) < v.min {
		resp.Error = validatorfuncerr.InvalidParameterValueFuncError(
			req.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		)
	}
}

// SizeAtLeast returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a Set.
//   - Contains at least min elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeAtLeast(minVal int) sizeAtLeastValidator {
	return sizeAtLeastValidator{
		min: minVal,
	}
}
