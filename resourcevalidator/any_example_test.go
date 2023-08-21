// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
)

func ExampleAny() {
	// Used inside a resource.Resource type ConfigValidators method
	_ = []resource.ConfigValidator{
		resourcevalidator.Any( /* ... */ ),
	}
}
