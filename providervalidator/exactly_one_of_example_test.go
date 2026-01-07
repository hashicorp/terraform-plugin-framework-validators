// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package providervalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func ExampleExactlyOneOf() {
	// Used inside a provider.Provider type ConfigValidators method
	_ = []provider.ConfigValidator{
		// Validate only one of the schema defined attributes named attr1
		// and attr2 has a known, non-null value.
		providervalidator.ExactlyOneOf(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
