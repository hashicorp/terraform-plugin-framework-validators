package providervalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func ExampleAtLeastOneOf() {
	// Used inside a provider.Provider type ConfigValidators method
	_ = []provider.ConfigValidator{
		// Validate at least one of the schema defined attributes named attr1
		// and attr2 has a known, non-null value.
		providervalidator.AtLeastOneOf(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
