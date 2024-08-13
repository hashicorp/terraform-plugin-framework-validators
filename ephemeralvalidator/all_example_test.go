// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ephemeralvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"

	"github.com/hashicorp/terraform-plugin-framework-validators/ephemeralvalidator"
)

func ExampleAll() {
	// Used inside a ephemeral.EphemeralResource type ConfigValidators method
	_ = []ephemeral.ConfigValidator{
		// The configuration must satisfy either All validator.
		ephemeralvalidator.Any(
			ephemeralvalidator.All( /* ... */ ),
			ephemeralvalidator.All( /* ... */ ),
		),
	}
}
