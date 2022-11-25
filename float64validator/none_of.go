package float64validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.Float64 = noneOfValidator{}

// noneOfValidator validates that the value does not match one of the values.
type noneOfValidator struct {
	values []types.Float64
}

func (v noneOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v noneOfValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be none of: %q", v.values)
}

func (v noneOfValidator) ValidateFloat64(ctx context.Context, request validator.Float64Request, response *validator.Float64Response) {
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

// NoneOf checks that the float64 held in the attribute
// is none of the given `values`.
func NoneOf(values ...float64) validator.Float64 {
	frameworkValues := make([]types.Float64, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.Float64Value(value))
	}

	return noneOfValidator{
		values: frameworkValues,
	}
}
