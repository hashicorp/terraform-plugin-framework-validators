package datasourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
)

func ExampleAny() {
	// Used inside a datasource.DataSource type ConfigValidators method
	_ = []datasource.ConfigValidator{
		datasourcevalidator.Any( /* ... */ ),
	}
}
