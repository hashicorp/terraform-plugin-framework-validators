// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package listresourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listresourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ExampleAtLeastOneOf() {
	// Used inside a list.ListResource type ConfigValidators method
	_ = []list.ConfigValidator{
		// Validate at least one of the schema defined attributes named attr1
		// and attr2 has a known, non-null value.
		listresourcevalidator.AtLeastOneOf(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
