package numbervalidator

import (
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/primitivevalidator"
)

// OneOf checks that the *big.Float held in the attribute
// is one of the given `acceptableFloats`.
func OneOf(acceptableFloats ...*big.Float) tfsdk.AttributeValidator {
	acceptableFloatValues := make([]attr.Value, 0, len(acceptableFloats))
	for _, f := range acceptableFloats {
		acceptableFloatValues = append(acceptableFloatValues, types.NumberValue(f))
	}

	return primitivevalidator.OneOf(acceptableFloatValues...)
}
