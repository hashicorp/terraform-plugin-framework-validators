package providervalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
)

func ExampleAny() {
	// Used inside a provider.Provider type ConfigValidators method
	_ = []provider.ConfigValidator{
		providervalidator.Any( /* ... */ ),
	}
}
