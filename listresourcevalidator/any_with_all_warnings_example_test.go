// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package listresourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listresourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
)

func ExampleAnyWithAllWarnings() {
	// Used inside a list.ListResource type ConfigValidators method
	_ = []list.ConfigValidator{
		listresourcevalidator.AnyWithAllWarnings( /* ... */ ),
	}
}
