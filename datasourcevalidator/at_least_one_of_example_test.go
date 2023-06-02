// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ExampleAtLeastOneOf() {
	// Used inside a datasource.DataSource type ConfigValidators method
	_ = []datasource.ConfigValidator{
		// Validate at least one of the schema defined attributes named attr1
		// and attr2 has a known, non-null value.
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
