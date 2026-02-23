// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package statestorevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/statestorevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

func ExampleConflicting() {
	// Used inside a statestore.StateStore type ConfigValidators method
	_ = []statestore.ConfigValidator{
		// Validate that schema defined attributes named attr1 and attr2 are not
		// both configured with known, non-null values.
		statestorevalidator.Conflicting(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
