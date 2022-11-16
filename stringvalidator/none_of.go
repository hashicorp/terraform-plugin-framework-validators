package stringvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/primitivevalidator"
)

// NoneOf checks that the string held in the attribute
// is none of the given `unacceptableStrings`.
func NoneOf(unacceptableStrings ...string) tfsdk.AttributeValidator {
	unacceptableStringValues := make([]attr.Value, 0, len(unacceptableStrings))
	for _, s := range unacceptableStrings {
		unacceptableStringValues = append(unacceptableStringValues, types.StringValue(s))
	}

	return primitivevalidator.NoneOf(unacceptableStringValues...)
}

// NoneOfCaseInsensitive checks that the string held in the attribute
// is none of the given `unacceptableStrings`, irrespective of case sensitivity.
func NoneOfCaseInsensitive(unacceptableStrings ...string) tfsdk.AttributeValidator {
	return &acceptableStringsCaseInsensitiveAttributeValidator{
		unacceptableStrings,
		false,
	}
}
