package validatordiag

import (
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// InvalidAttributeValueDiagnostic returns an error Diagnostic to be used when an attribute has an invalid value.
func InvalidAttributeValueDiagnostic(path path.Path, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value",
		capitalize(description)+", got: "+value,
	)
}

// InvalidAttributeValueLengthDiagnostic returns an error Diagnostic to be used when an attribute's value has an invalid length.
func InvalidAttributeValueLengthDiagnostic(path path.Path, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value Length",
		capitalize(description)+", got: "+value,
	)
}

// InvalidAttributeValueMatchDiagnostic returns an error Diagnostic to be used when an attribute's value has an invalid match.
func InvalidAttributeValueMatchDiagnostic(path path.Path, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value Match",
		capitalize(description)+", got: "+value,
	)
}

// InvalidAttributeCombinationDiagnostic returns an error Diagnostic to be used when a schemavalidator of attributes is invalid.
func InvalidAttributeCombinationDiagnostic(path path.Path, description string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Combination",
		capitalize(description),
	)
}

// InvalidAttributeTypeDiagnostic returns an error Diagnostic to be used when an attribute has an invalid type.
func InvalidAttributeTypeDiagnostic(path path.Path, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Type",
		capitalize(description)+", got: "+value,
	)
}

func BugInProviderDiagnostic(summary string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(summary,
		"This is a bug in the provider, which should be reported in the provider's own issue tracker",
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
