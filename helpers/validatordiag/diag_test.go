package validatordiag

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
			got := capitalize(test.input)

			if got != test.expected {
				t.Fatalf("expected: %q, got: %q", test.expected, got)
			}
		})
	}
}

func TestGetErrors(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input    diag.Diagnostics
		expected diag.Diagnostics
	}
	tests := map[string]testCase{
		"errors": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error", ""),
				diag.NewWarningDiagnostic("Warning", ""),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error", ""),
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			got := GetErrors(test.input)

			if diff := cmp.Diff(test.expected, got); diff != "" {
				t.Fatalf("expected: %q, got: %q", test.expected, got)
			}
		})
	}
}

func TestGetWarnings(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input    diag.Diagnostics
		expected diag.Diagnostics
	}
	tests := map[string]testCase{
		"errors": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error", ""),
				diag.NewWarningDiagnostic("Warning", ""),
			},
			expected: diag.Diagnostics{
				diag.NewWarningDiagnostic("Warning", ""),
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			got := GetWarnings(test.input)

			if diff := cmp.Diff(test.expected, got); diff != "" {
				t.Fatalf("expected: %q, got: %q", test.expected, got)
			}
		})
	}
}
