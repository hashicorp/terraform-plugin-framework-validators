package float64validator

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/primitivevalidator"
)

// NoneOf checks that the float64 held in the attribute
// is none of the given `unacceptableFloats`.
func NoneOf(unacceptableFloats ...float64) tfsdk.AttributeValidator {
	unacceptableFloatValues := make([]attr.Value, 0, len(unacceptableFloats))
	for _, f := range unacceptableFloats {
		unacceptableFloatValues = append(unacceptableFloatValues, types.Float64Value(f))
	}

	return primitivevalidator.NoneOf(unacceptableFloatValues...)
}
