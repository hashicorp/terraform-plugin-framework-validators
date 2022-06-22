package validatordiag

import (
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// InvalidValueDiagnostic returns an error Diagnostic to be used when an attribute has an invalid value.
func InvalidValueDiagnostic(path *tftypes.AttributePath, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value",
		capitalize(description)+", got: "+value,
	)
}

// InvalidValueLengthDiagnostic returns an error Diagnostic to be used when an attribute's value has an invalid length.
func InvalidValueLengthDiagnostic(path *tftypes.AttributePath, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value Length",
		capitalize(description)+", got: "+value,
	)
}

// InvalidValueMatchDiagnostic returns an error Diagnostic to be used when an attribute's value has an invalid match.
func InvalidValueMatchDiagnostic(path *tftypes.AttributePath, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value Match",
		capitalize(description)+", got: "+value,
	)
}

// InvalidSchemaDiagnostic returns an error Diagnostic to be used when a schemavalidator of attributes is invalid.
func InvalidSchemaDiagnostic(path *tftypes.AttributePath, description string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Combination",
		capitalize(description),
	)
}

// InvalidTypeDiagnostic returns an error Diagnostic to be used when an attribute has an invalid type.
func InvalidTypeDiagnostic(path *tftypes.AttributePath, description string, value string) diag.Diagnostic {
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
