// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.Float64 = atLeastValidator{}

// atLeastValidator validates that an float Attribute's value is at least a certain value.
type atLeastValidator struct {
	min float64
}

// Description describes the validation in plain text formatting.
func (validator atLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at least %f", validator.min)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator atLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// ValidateFloat64 performs the validation.
func (validator atLeastValidator) ValidateFloat64(ctx context.Context, request validator.Float64Request, response *validator.Float64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueFloat64()

	if value < validator.min {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			validator.Description(ctx),
			fmt.Sprintf("%f", value),
		))
	}
}

// AtLeast returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a number, which can be represented by a 64-bit floating point.
//   - Is greater than or equal to the given minimum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtLeast(min float64) validator.Float64 {
	return atLeastValidator{
		min: min,
	}
}
