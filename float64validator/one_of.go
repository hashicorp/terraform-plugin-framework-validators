package float64validator

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/primitivevalidator"
)

// OneOf checks that the float64 held in the attribute
// is one of the given `acceptableFloats`.
func OneOf(acceptableFloats ...float64) tfsdk.AttributeValidator {
	acceptableFloatValues := make([]attr.Value, 0, len(acceptableFloats))
	for _, f := range acceptableFloats {
		acceptableFloatValues = append(acceptableFloatValues, types.Float64Value(f))
	}

	return primitivevalidator.OneOf(acceptableFloatValues...)
}
