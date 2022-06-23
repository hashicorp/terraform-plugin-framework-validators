package stringvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// OneOf checks that the string held in the attribute
// is one of the given `acceptableStrings`.
//
// String comparison case sensitiveness is controlled by the `caseSensitive` argument.
func OneOf(caseSensitive bool, acceptableStrings ...string) tfsdk.AttributeValidator {
	return &acceptableStringsAttributeValidator{
		acceptableStrings,
		caseSensitive,
		true,
	}
}
