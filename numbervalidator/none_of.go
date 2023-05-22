// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package numbervalidator

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.Number = noneOfValidator{}

// noneOfValidator validates that the value does not match one of the values.
type noneOfValidator struct {
	values []types.Number
}

func (v noneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v noneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be none of: %q", v.values)
}

func (v noneOfValidator) ValidateNumber(ctx context.Context, request validator.NumberRequest, response *validator.NumberResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, otherValue := range v.values {
		if !value.Equal(otherValue) {
			continue
		}

		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value.String(),
		))

		break
	}
}

// NoneOf checks that the Number held in the attribute
// is none of the given `values`.
func NoneOf(values ...*big.Float) validator.Number {
	frameworkValues := make([]types.Number, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.NumberValue(value))
	}

	return noneOfValidator{
		values: frameworkValues,
	}
}
