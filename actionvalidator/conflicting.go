// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package actionvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// Conflicting checks that a set of path.Expression, are not configured
// simultaneously.
func Conflicting(expressions ...path.Expression) action.ConfigValidator {
	return &configvalidator.ConflictingValidator{
		PathExpressions: expressions,
	}
}
