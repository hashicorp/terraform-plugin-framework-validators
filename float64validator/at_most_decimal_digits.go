// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package float64validator

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Float64 = atMostDecimalDigitsValidator{}

// atMostDecimalDigitsValidator validates that an float Attribute's value has at most a certain number of decimal digits.
type atMostDecimalDigitsValidator struct {
	atMostDecimalDigits int
}

// Description describes the validation in plain text formatting.
func (validator atMostDecimalDigitsValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must have up to %d decimal digits", validator.atMostDecimalDigits)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator atMostDecimalDigitsValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// ValidateFloat64 performs the validation.
func (validator atMostDecimalDigitsValidator) ValidateFloat64(ctx context.Context, request validator.Float64Request, response *validator.Float64Response) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueFloat64()

	if decDigits := countDecimalDigits(value); decDigits > validator.atMostDecimalDigits {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			validator.Description(ctx),
			strconv.Itoa(decDigits),
		))
	}
}

// countDecimalDigits returns the number of decimal digits in a float64.
func countDecimalDigits(f float64) int {
	str := strconv.FormatFloat(f, 'f', -1, 64)
	str = strings.TrimRight(str, "0")
	parts := strings.Split(str, ".")
	if len(parts) != 2 {
		return 0
	}
	return len(parts[1])
}

// AtMostDecimalDigits returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - It contains less than or equal to the given maximum number of decimal digits.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func AtMostDecimalDigits(atMostDecimalDigits int) validator.Float64 {
	if atMostDecimalDigits < 0 {
		return nil
	}
	return atMostDecimalDigitsValidator{
		atMostDecimalDigits: atMostDecimalDigits,
	}
}
