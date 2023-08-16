package datasourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
)

func ExampleAny() {
	// Used inside a datasource.DataSource type ConfigValidators method
	_ = []datasource.ConfigValidator{
		// Validate that the configuration has either:
		// 	- only one of the schema defined attributes named attr1
		//  and attr2 has a known, non-null value OR
		//  - only one of the schema defined attributes named attr3
		//  and attr4 has a known, non-null value
		datasourcevalidator.Any(
			datasourcevalidator.ExactlyOneOf(
				path.MatchRoot("attr1"),
				path.MatchRoot("attr2"),
			),
			datasourcevalidator.ExactlyOneOf(
				path.MatchRoot("attr3"),
				path.MatchRoot("attr4"),
			),
		),
	}
}
