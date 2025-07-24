// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listresourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listresourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
)

func ExampleAll() {
	// Used inside a list.ListResource type ConfigValidators method
	_ = []list.ConfigValidator{
		// The configuration must satisfy either All validator.
		listresourcevalidator.Any(
			listresourcevalidator.All( /* ... */ ),
			listresourcevalidator.All( /* ... */ ),
		),
	}
}
