// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actionvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/actionvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ExampleRequiredTogether() {
	// Used inside a action.Action type ConfigValidators method
	_ = []action.ConfigValidator{
		// Validate the schema defined attributes named attr1 and attr2 are either
		// both null or both known values.
		actionvalidator.RequiredTogether(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
