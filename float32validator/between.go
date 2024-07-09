// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float32validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.Float32 = betweenValidator{}

// betweenValidator validates that an float Attribute's value is in a range.
type betweenValidator struct {
	min, max float32
}

// Description describes the validation in plain text formatting.
func (validator betweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be between %f and %f", validator.min, validator.max)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator betweenValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// ValidateFloat32 performs the validation.
func (v betweenValidator) ValidateFloat32(ctx context.Context, request validator.Float32Request, response *validator.Float32Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueFloat32()

	if value < v.min || value > v.max {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%f", value),
		))
	}
}

// Between returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a number, which can be represented by a 32-bit floating point.
//   - Is greater than or equal to the given minimum and less than or equal to the given maximum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func Between(min, max float32) validator.Float32 {
	if min > max {
		return nil
	}

	return betweenValidator{
		min: min,
		max: max,
	}
}
