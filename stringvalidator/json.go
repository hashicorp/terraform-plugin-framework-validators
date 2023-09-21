// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = IsJsonValidator{}

// IsJsonValidator validates that a string is valid Json.
type IsJsonValidator struct{}

// Description describes the validation in plain text formatting.
func (validator IsJsonValidator) Description(_ context.Context) string {
	return "string must be valid json"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator IsJsonValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func NormalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	bytes, _ := json.Marshal(j)
	return string(bytes[:]), nil
}

// Validate performs the validation.
func (v IsJsonValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if _, err := NormalizeJsonString(value); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%s", v),
		))

		return
	}
}

// IsJson returns an validator which validates string value is valid json
func IsJson() validator.String {
	return IsJsonValidator{}
}
