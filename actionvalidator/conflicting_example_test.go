// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actionvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/actionvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ExampleConflicting() {
	// Used inside a action.Action type ConfigValidators method
	_ = []action.ConfigValidator{
		// Validate that schema defined attributes named attr1 and attr2 are not
		// both configured with known, non-null values.
		actionvalidator.Conflicting(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
