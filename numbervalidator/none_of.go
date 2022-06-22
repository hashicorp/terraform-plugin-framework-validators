package numbervalidator

import (
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/primitivevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NoneOf checks that the *big.Float held in the attribute
// is none of the given `unacceptableFloats`.
func NoneOf(unacceptableFloats ...*big.Float) tfsdk.AttributeValidator {
	unacceptableFloatValues := make([]attr.Value, 0, len(unacceptableFloats))
	for _, f := range unacceptableFloats {
		unacceptableFloatValues = append(unacceptableFloatValues, types.Number{Value: f})
	}

	return primitivevalidator.NoneOf(unacceptableFloatValues...)
}
