// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package actionvalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/action"

	"github.com/hashicorp/terraform-plugin-framework-validators/actionvalidator"
)

func ExampleAll() {
	// Used inside a action.Action type ConfigValidators method
	_ = []action.ConfigValidator{
		// The configuration must satisfy either All validator.
		actionvalidator.Any(
			actionvalidator.All( /* ... */ ),
			actionvalidator.All( /* ... */ ),
		),
	}
}
