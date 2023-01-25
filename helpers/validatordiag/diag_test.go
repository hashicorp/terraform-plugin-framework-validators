package validatordiag

import (
	"testing"
)

func TestCapitalize(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input    string
		expected string
	}
	tests := map[string]testCase{
		"empty string": {
			input:    "",
			expected: "",
		},
		"all lowercase": {
			input:    "abcd",
			expected: "Abcd",
		},
		"all uppercase": {
			input:    "ABCD",
			expected: "ABCD",
		},
		"initial numeric": {
			input:    "1 ab",
			expected: "1 ab",
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := capitalize(test.input)

			if got != test.expected {
				t.Fatalf("expected: %q, got: %q", test.expected, got)
			}
		})
	}
}
