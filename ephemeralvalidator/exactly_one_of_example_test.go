// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ephemeralvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/ephemeralvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ExampleExactlyOneOf() {
	// Used inside a ephemeral.EphemeralResource type ConfigValidators method
	_ = []ephemeral.ConfigValidator{
		// Validate only one of the schema defined attributes named attr1
		// and attr2 has a known, non-null value.
		ephemeralvalidator.ExactlyOneOf(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
