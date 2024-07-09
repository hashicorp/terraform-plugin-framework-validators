// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float32validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.Float32 = atLeastValidator{}

// atLeastValidator validates that an float Attribute's value is at least a certain value.
type atLeastValidator struct {
	min float32
}

// Description describes the validation in plain text formatting.
func (validator atLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at least %f", validator.min)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator atLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// ValidateFloat32 performs the validation.
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

// AtLeast returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a number, which can be represented by a 32-bit floating point.
//   - Is greater than or equal to the given minimum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtLeast(min float32) validator.Float32 {
	return atLeastValidator{
		min: min,
	}
}
