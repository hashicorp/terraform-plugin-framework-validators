package int64validator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/primitivevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NoneOf checks that the int64 held in the attribute
// is none of the given `unacceptableInts`.
func NoneOf(unacceptableInts ...int64) tfsdk.AttributeValidator {
	unacceptableIntValues := make([]attr.Value, 0, len(unacceptableInts))
	for _, i := range unacceptableInts {
		unacceptableIntValues = append(unacceptableIntValues, types.Int64{Value: i})
	}

	return primitivevalidator.NoneOf(unacceptableIntValues...)
}
