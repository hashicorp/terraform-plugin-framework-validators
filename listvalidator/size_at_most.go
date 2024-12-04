// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatorfuncerr"
)

var _ validator.List = sizeAtMostValidator{}
var _ function.ListParameterValidator = sizeAtMostValidator{}

type sizeAtMostValidator struct {
	max int
}

func (v sizeAtMostValidator) Description(_ context.Context) string {
	return fmt.Sprintf("list must contain at most %d elements", v.max)
}

func (v sizeAtMostValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v sizeAtMostValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
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
		return
	}

	elems := req.ConfigValue.Elements()

	if len(elems) > v.max {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			req.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		))
	}
}

func (v sizeAtMostValidator) ValidateParameterList(ctx context.Context, req function.ListParameterValidatorRequest, resp *function.ListParameterValidatorResponse) {
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
		return
	}

	elems := req.Value.Elements()

	if len(elems) > v.max {
		resp.Error = validatorfuncerr.InvalidParameterValueFuncError(
			req.ArgumentPosition,
			v.Description(ctx),
			fmt.Sprintf("%d", len(elems)),
		)
	}
}

// SizeAtMost returns an AttributeValidator which ensures that any configured
// attribute or function parameter value:
//
//   - Is a List.
//   - Contains at most max elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeAtMost(maxVal int) sizeAtMostValidator {
	return sizeAtMostValidator{
		max: maxVal,
	}
}
