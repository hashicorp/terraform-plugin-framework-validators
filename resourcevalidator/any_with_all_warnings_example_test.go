package resourcevalidator_test

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
)

func ExampleAnyWithAllWarnings() {
	// Used inside a resource.Resource type ConfigValidators method
	_ = []resource.ConfigValidator{
		resourcevalidator.AnyWithAllWarnings( /* ... */ ),
	}
}
