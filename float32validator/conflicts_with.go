// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package float32validator

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator"
)

// ConflictsWith checks that a set of path.Expression,
// including the attribute the validator is applied to,
// do not have a value simultaneously.
//
// This implements the validation logic declaratively within the schema.
// Refer to [datasourcevalidator.Conflicting],
// [providervalidator.Conflicting], or [resourcevalidator.Conflicting]
// for declaring this type of validation outside the schema definition.
//
// Relative path.Expression will be resolved using the attribute being
// validated.
func ConflictsWith(expressions ...path.Expression) validator.Float32 {
	return schemavalidator.ConflictsWithValidator{
		PathExpressions: expressions,
	}
}
