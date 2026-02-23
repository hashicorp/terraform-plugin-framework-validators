// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package statestorevalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

// AtLeastOneOf checks that a set of path.Expression has at least one non-null
// or unknown value.
func AtLeastOneOf(expressions ...path.Expression) statestore.ConfigValidator {
	return &configvalidator.AtLeastOneOfValidator{
		PathExpressions: expressions,
	}
}
