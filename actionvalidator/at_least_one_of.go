// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actionvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// AtLeastOneOf checks that a set of path.Expression has at least one non-null
// or unknown value.
func AtLeastOneOf(expressions ...path.Expression) action.ConfigValidator {
	return &configvalidator.AtLeastOneOfValidator{
		PathExpressions: expressions,
	}
}
