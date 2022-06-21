package validatordiag

import (
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// AttributeValueDiagnostic returns an error Diagnostic to be used when an attribute has an invalid value.
func AttributeValueDiagnostic(path *tftypes.AttributePath, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value",
		capitalize(description)+", got: "+value,
	)
}

// AttributeValueLengthDiagnostic returns an error Diagnostic to be used when an attribute's value has an invalid length.
func AttributeValueLengthDiagnostic(path *tftypes.AttributePath, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value Length",
		capitalize(description)+", got: "+value,
	)
}

// AttributeValueMatchesDiagnostic returns an error Diagnostic to be used when an attribute's value has an invalid match.
func AttributeValueMatchesDiagnostic(path *tftypes.AttributePath, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value Match",
		capitalize(description)+", got: "+value,
	)
}

// AttributeValueTerraformValueDiagnostic returns an error Diagnostic to be used when an attribute's value cannot be
// converted to terraform value.
func AttributeValueTerraformValueDiagnostic(path *tftypes.AttributePath, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value Terraform Value",
		capitalize(description)+", err: "+value,
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
