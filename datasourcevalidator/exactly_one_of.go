package datasourcevalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// ExactlyOneOf checks that a set of path.Expression does not have more than
// one known value.
func ExactlyOneOf(expressions ...path.Expression) datasource.ConfigValidator {
	return &configvalidator.ExactlyOneOfValidator{
		PathExpressions: expressions,
	}
}
