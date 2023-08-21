package stringvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// NotEqual checks that a set of path.Expression have values not equal to
// the current attribute when the current attribute is non-null.
//
// Relative path.Expression will be resolved using the attribute being
// validated.
func NotEqual(expressions ...path.Expression) validator.String {
	return &schemavalidator.NotEqualValidator{
		PathExpressions: expressions,
	}
}
