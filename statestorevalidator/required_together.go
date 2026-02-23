// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package statestorevalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

// RequiredTogether checks that a set of path.Expression either has all known
// or all null values.
func RequiredTogether(expressions ...path.Expression) statestore.ConfigValidator {
	return &configvalidator.RequiredTogetherValidator{
		PathExpressions: expressions,
	}
}
