package int64validator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/primitivevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OneOf checks that the int64 held in the attribute
// is one of the given `acceptableInts`.
func OneOf(acceptableInts ...int64) tfsdk.AttributeValidator {
	acceptableIntValues := make([]attr.Value, 0, len(acceptableInts))
	for _, i := range acceptableInts {
		acceptableIntValues = append(acceptableIntValues, types.Int64{Value: i})
	}

	return primitivevalidator.OneOf(acceptableIntValues...)
}
