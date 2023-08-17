package providervalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
)

func ExampleAll() {
	// Used inside a provider.Provider type ConfigValidators method
	_ = []provider.ConfigValidator{
		// The configuration must satisfy either All validator.
		providervalidator.Any(
			providervalidator.All( /* ... */ ),
			providervalidator.All( /* ... */ ),
		),
	}
}
