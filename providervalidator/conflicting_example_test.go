package providervalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func ExampleConflicting() {
	// Used inside a provider.Provider type ConfigValidators method
	_ = []provider.ConfigValidator{
		// Validate that schema defined attributes named attr1 and attr2 are not
		// both configured with known, non-null values.
		providervalidator.Conflicting(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
