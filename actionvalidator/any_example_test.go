// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package actionvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/actionvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
)

func ExampleAny() {
	// Used inside a action.Action type ConfigValidators method
	_ = []action.ConfigValidator{
		actionvalidator.Any( /* ... */ ),
	}
}
