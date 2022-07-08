package pathutils

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// MergeExpressionsWithAttribute returns the given path.Expressions,
// but each has been merged with the given attribute path.Expression,
// and then resolved.
//
// Additionally, if the attribute path.Expression was not part of the initial slice,
// it is added to the result.
func MergeExpressionsWithAttribute(pathExps path.Expressions, attrPathExp path.Expression) path.Expressions {
	result := make(path.Expressions, 0, len(pathExps)+1)

	// First, add the attribute own path expression to the result
	result.Append(attrPathExp)
	// Then, add all the other path expressions,
	// after they have been merged to the attribute own path
	for _, pe := range pathExps {
		result.Append(attrPathExp.Merge(pe))
	}

	return result
}
