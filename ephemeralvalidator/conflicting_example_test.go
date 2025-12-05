// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package ephemeralvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/ephemeralvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func ExampleConflicting() {
	// Used inside a ephemeral.EphemeralResource type ConfigValidators method
	_ = []ephemeral.ConfigValidator{
		// Validate that schema defined attributes named attr1 and attr2 are not
		// both configured with known, non-null values.
		ephemeralvalidator.Conflicting(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
