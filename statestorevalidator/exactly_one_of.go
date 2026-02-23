// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package statestorevalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

// ExactlyOneOf checks that a set of path.Expression does not have more than
// one known value.
func ExactlyOneOf(expressions ...path.Expression) statestore.ConfigValidator {
	return &configvalidator.ExactlyOneOfValidator{
		PathExpressions: expressions,
	}
}
