package numbervalidator

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.Number = oneOfValidator{}

// oneOfValidator validates that the value matches one of expected values.
type oneOfValidator struct {
	values []types.Number
}

func (v oneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v oneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be one of: %q", v.values)
}

func (v oneOfValidator) ValidateNumber(ctx context.Context, request validator.NumberRequest, response *validator.NumberResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	for _, otherValue := range v.values {
		if value.Equal(otherValue) {
			return
		}
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
		request.Path,
		v.Description(ctx),
		value.String(),
	))
}

// OneOf checks that the Number held in the attribute
// is none of the given `values`.
func OneOf(values ...*big.Float) validator.Number {
	frameworkValues := make([]types.Number, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.NumberValue(value))
	}

	return oneOfValidator{
		values: frameworkValues,
	}
}
