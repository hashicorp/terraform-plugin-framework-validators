// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package statestorevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/statestore"

	"github.com/hashicorp/terraform-plugin-framework-validators/statestorevalidator"
)

func ExampleAll() {
	// Used inside a statestore.StateStore type ConfigValidators method
	_ = []statestore.ConfigValidator{
		// The configuration must satisfy either All validator.
		statestorevalidator.Any(
			statestorevalidator.All( /* ... */ ),
			statestorevalidator.All( /* ... */ ),
		),
	}
}
