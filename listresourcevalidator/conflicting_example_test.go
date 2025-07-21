// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listresourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listresourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ExampleConflicting() {
	// Used inside a list.ListResource type ConfigValidators method
	_ = []list.ConfigValidator{
		// Validate that schema defined attributes named attr1 and attr2 are not
		// both configured with known, non-null values.
		listresourcevalidator.Conflicting(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
