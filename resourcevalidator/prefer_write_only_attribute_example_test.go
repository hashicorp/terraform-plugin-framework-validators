// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
)

func ExamplePreferWriteOnlyAttribute() {
	// Used inside a resource.Resource type ConfigValidators method
	_ = []resource.ConfigValidator{
		// Throws a warning diagnostic if the resource supports write-only
		// attributes and the oldAttribute has a known value.
		resourcevalidator.PreferWriteOnlyAttribute(
			path.MatchRoot("oldAttribute"),
			path.MatchRoot("writeOnlyAttribute"),
		),
	}
}
