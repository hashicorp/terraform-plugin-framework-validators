// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package datasourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ExampleRequiredTogether() {
	// Used inside a datasource.DataSource type ConfigValidators method
	_ = []datasource.ConfigValidator{
		// Validate the schema defined attributes named attr1 and attr2 are either
		// both null or both known values.
		datasourcevalidator.RequiredTogether(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
