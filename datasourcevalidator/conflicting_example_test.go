// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ExampleConflicting() {
	// Used inside a datasource.DataSource type ConfigValidators method
	_ = []datasource.ConfigValidator{
		// Validate that schema defined attributes named attr1 and attr2 are not
		// both configured with known, non-null values.
		datasourcevalidator.Conflicting(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
