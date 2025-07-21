// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listresourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listresourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ExampleExactlyOneOf() {
	// Used inside a list.ListResource type ConfigValidators method
	_ = []list.ConfigValidator{
		// Validate only one of the schema defined attributes named attr1
		// and attr2 has a known, non-null value.
		listresourcevalidator.ExactlyOneOf(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
