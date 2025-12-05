// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package actionvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/actionvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ExampleExactlyOneOf() {
	// Used inside a action.Action type ConfigValidators method
	_ = []action.ConfigValidator{
		// Validate only one of the schema defined attributes named attr1
		// and attr2 has a known, non-null value.
		actionvalidator.ExactlyOneOf(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
