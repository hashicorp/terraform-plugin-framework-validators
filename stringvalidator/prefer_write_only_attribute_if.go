// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package stringvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator"
)

// PreferWriteOnlyAttributeIfFunc is a conditional function used in the
// PreferWriteOnlyAttributeIf validator to determine whether a warning should be returned.
type PreferWriteOnlyAttributeIfFunc func(context.Context, validator.StringRequest, *PreferWriteOnlyAttributeIfFuncResponse)

// PreferWriteOnlyAttributeIfFuncResponse is the response type for a
// PreferWriteOnlyAttributeIfFunc.
type PreferWriteOnlyAttributeIfFuncResponse = schemavalidator.PreferWriteOnlyAttributeIfFuncResponse

// PreferWriteOnlyAttributeIf returns a warning if the given function enables it
// and the PreferWriteOnlyAttribute validation conditions are met.
func PreferWriteOnlyAttributeIf(writeOnlyAttribute path.Expression, f PreferWriteOnlyAttributeIfFunc) validator.String {
	return schemavalidator.PreferWriteOnlyAttributeIf{
		WriteOnlyAttribute: writeOnlyAttribute,
		StringFunc:         f,
	}
}
