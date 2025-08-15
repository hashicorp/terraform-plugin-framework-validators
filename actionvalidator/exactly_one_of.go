// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actionvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// ExactlyOneOf checks that a set of path.Expression does not have more than
// one known value.
func ExactlyOneOf(expressions ...path.Expression) action.ConfigValidator {
	return &configvalidator.ExactlyOneOfValidator{
		PathExpressions: expressions,
	}
}
