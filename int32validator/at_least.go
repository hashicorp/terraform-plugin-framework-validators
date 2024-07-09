// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package int32validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.Int32 = atLeastValidator{}

// atLeastValidator validates that an integer Attribute's value is at least a certain value.
type atLeastValidator struct {
	min int32
}

// Description describes the validation in plain text formatting.
func (validator atLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be at least %d", validator.min)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator atLeastValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// ValidateInt32 performs the validation.
func (v atLeastValidator) ValidateInt32(ctx context.Context, request validator.Int32Request, response *validator.Int32Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if request.ConfigValue.ValueInt32() < v.min {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", request.ConfigValue.ValueInt32()),
		))
	}
}

// AtLeast returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a number, which can be represented by a 32-bit integer.
//   - Is greater than or equal to the given minimum.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtLeast(min int32) validator.Int32 {
	return atLeastValidator{
		min: min,
	}
}
