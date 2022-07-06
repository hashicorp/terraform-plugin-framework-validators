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
	result = append(result, attrPathExp)

	for _, pe := range pathExps {
		mpe := attrPathExp.Merge(pe).Resolve()

		// Include the merged path expression,
		// only if it's not the same as the attribute
		if !mpe.Equal(attrPathExp) {
			result = append(result, mpe)
		}
	}

	return result
}
