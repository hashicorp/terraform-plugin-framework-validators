package validatordiag

import (
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// AttributeValueDiagnostic returns an error Diagnostic to be used when an attribute has an unexpected value.
func AttributeValueDiagnostic(path *tftypes.AttributePath, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value",
		capitalize(description)+", got: "+value,
	)
}

// capitalize will uppercase the first letter in a UTF-8 string.
func capitalize(str string) string {
	if str == "" {
		return ""
	}

	firstRune, size := utf8.DecodeRuneInString(str)

	return string(unicode.ToUpper(firstRune)) + str[size:]
}
