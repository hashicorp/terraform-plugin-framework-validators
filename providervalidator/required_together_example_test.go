package providervalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func ExampleRequiredTogether() {
	// Used inside a provider.Provider type ConfigValidators method
	_ = []provider.ConfigValidator{
		// Validate the schema defined attributes named attr1 and attr2 are either
		// both null or both known values.
		providervalidator.RequiredTogether(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
