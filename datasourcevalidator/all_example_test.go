// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
)

func ExampleAll() {
	// Used inside a datasource.DataSource type ConfigValidators method
	_ = []datasource.ConfigValidator{
		// The configuration must satisfy either All validator.
		datasourcevalidator.Any(
			datasourcevalidator.All( /* ... */ ),
			datasourcevalidator.All( /* ... */ ),
		),
	}
}
