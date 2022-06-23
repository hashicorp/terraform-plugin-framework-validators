package stringvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// NoneOf checks that the string held in the attribute
// is none of the given `unacceptableStrings`.
//
// String comparison case sensitiveness is controlled by the `caseSensitive` argument.
func NoneOf(caseSensitive bool, unacceptableStrings ...string) tfsdk.AttributeValidator {
	return &acceptableStringsAttributeValidator{
		unacceptableStrings,
		caseSensitive,
		false,
	}
}
