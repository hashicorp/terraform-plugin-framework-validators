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

// InvalidAttributeSchemaDiagnostic returns an error Diagnostic to be used when a schemavalidator of attributes is invalid.
func InvalidAttributeSchemaDiagnostic(path path.Path, description string) diag.Diagnostic {
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

// ErrorsCount returns the amount of diag.Diagnostic in diag.Diagnostics that are diag.SeverityError.
func ErrorsCount(diags diag.Diagnostics) int {
	count := 0

	for _, d := range diags {
		if diag.SeverityError == d.Severity() {
			count++
		}
	}

	return count
}

// WarningsCount returns the amount of diag.Diagnostic in diag.Diagnostics that are diag.SeverityWarning.
func WarningsCount(diags diag.Diagnostics) int {
	count := 0

	for _, d := range diags {
		if diag.SeverityWarning == d.Severity() {
			count++
		}
	}

	return count
}

// capitalize will uppercase the first letter in a UTF-8 string.
func capitalize(str string) string {
	if str == "" {
		return ""
	}

	firstRune, size := utf8.DecodeRuneInString(str)

	return string(unicode.ToUpper(firstRune)) + str[size:]
}
