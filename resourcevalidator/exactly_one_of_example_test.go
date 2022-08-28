package resourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func ExampleExactlyOneOf() {
	// Used inside a resource.Resource type ConfigValidators method
	_ = []resource.ConfigValidator{
		// Validate only one of the schema defined attributes named attr1
		// and attr2 has a known, non-null value.
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("attr1"),
			path.MatchRoot("attr2"),
		),
	}
}
