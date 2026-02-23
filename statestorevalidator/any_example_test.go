// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package statestorevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/statestorevalidator"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

func ExampleAny() {
	// Used inside a statestore.StateStore type ConfigValidators method
	_ = []statestore.ConfigValidator{
		statestorevalidator.Any( /* ... */ ),
	}
}
