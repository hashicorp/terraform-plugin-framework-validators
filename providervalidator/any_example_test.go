// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

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
