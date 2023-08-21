// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package providervalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
)

func ExampleAnyWithAllWarnings() {
	// Used inside a provider.Provider type ConfigValidators method
	_ = []provider.ConfigValidator{
		providervalidator.AnyWithAllWarnings( /* ... */ ),
	}
}
