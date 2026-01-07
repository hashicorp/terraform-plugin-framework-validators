// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package ephemeralvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"

	"github.com/hashicorp/terraform-plugin-framework-validators/ephemeralvalidator"
)

func ExampleAnyWithAllWarnings() {
	// Used inside a ephemeral.EphemeralResource type ConfigValidators method
	_ = []ephemeral.ConfigValidator{
		ephemeralvalidator.AnyWithAllWarnings( /* ... */ ),
	}
}
