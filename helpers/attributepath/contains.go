package attributepath

import (
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Contains returns true if needle (one *tftypes.AttributePath)
// can be found in haystack (collection of *tftypes.AttributePath).
func Contains(needle *tftypes.AttributePath, haystack ...*tftypes.AttributePath) bool {
	for _, p := range haystack {
		if needle.Equal(p) {
			return true
		}
	}
	return false
}
