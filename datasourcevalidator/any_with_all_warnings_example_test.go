// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package datasourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
)

func ExampleAnyWithAllWarnings() {
	// Used inside a datasource.DataSource type ConfigValidators method
	_ = []datasource.ConfigValidator{
		datasourcevalidator.AnyWithAllWarnings( /* ... */ ),
	}
}
