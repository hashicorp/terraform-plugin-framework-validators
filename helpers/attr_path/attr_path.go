package attr_path

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ToString takes all the tftypes.AttributePathStep in a tftypes.AttributePath and concatenates them,
// using `.` as separator.
//
// This should be used only when trying to "print out" a tftypes.AttributePath in a log or an error message.
func ToString(path *tftypes.AttributePath) string {
	var res strings.Builder
	for pos, step := range path.Steps() {
		if pos != 0 {
			res.WriteString(".")
		}
		switch v := step.(type) {
		case tftypes.AttributeName:
			res.WriteString(string(v))
		case tftypes.ElementKeyString:
			res.WriteString(string(v))
		case tftypes.ElementKeyInt:
			res.WriteString(strconv.FormatInt(int64(v), 10))
		case tftypes.ElementKeyValue:
			res.WriteString(tftypes.Value(v).String())
		}
	}

	return res.String()
}

// JoinToString works similarly to strings.Join: it takes a collection of *tftypes.AttributePath,
// applies to each ToString, and the resulting strings with a `,` separator.
//
// This should be used only when trying to "print out" a tftypes.AttributePath in a log or an error message.
func JoinToString(paths ...*tftypes.AttributePath) string {
	res := make([]string, len(paths))
	for i, path := range paths {
		res[i] = ToString(path)
	}

	return strings.Join(res, ",")
}

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
