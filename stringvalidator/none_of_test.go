// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package stringvalidator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func TestNoneOfValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in           types.String
		noneOfValues []string
		expectError  bool
	}

	testCases := map[string]testCase{
		"simple-match": {
			in: types.StringValue("foo"),
			noneOfValues: []string{
				"foo",
				"bar",
				"baz",
			},
			expectError: true,
		},
		"simple-mismatch-case-insensitive": {
			in: types.StringValue("foo"),
			noneOfValues: []string{
				"FOO",
				"bar",
				"baz",
			},
		},
		"simple-mismatch": {
			in: types.StringValue("foz"),
			noneOfValues: []string{
				"foo",
				"bar",
				"baz",
			},
		},
		"skip-validation-on-null": {
			in: types.StringNull(),
			noneOfValues: []string{
				"foo",
				"bar",
				"baz",
			},
		},
		"skip-validation-on-unknown": {
			in: types.StringUnknown(),
			noneOfValues: []string{
				"foo",
				"bar",
				"baz",
			},
		},
	}

	for name, test := range testCases {

		t.Run(fmt.Sprintf("ValidateString - %s", name), func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				ConfigValue: test.in,
			}
			res := validator.StringResponse{}
			stringvalidator.NoneOf(test.noneOfValues...).ValidateString(context.TODO(), req, &res)

			if !res.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Diagnostics)
			}
		})

		t.Run(fmt.Sprintf("ValidateParameterString - %s", name), func(t *testing.T) {
			t.Parallel()
			req := function.StringParameterValidatorRequest{
				Value: test.in,
			}
			res := function.StringParameterValidatorResponse{}
			stringvalidator.NoneOf(test.noneOfValues...).ValidateParameterString(context.TODO(), req, &res)

			if res.Error == nil && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if res.Error != nil && !test.expectError {
				t.Fatalf("got unexpected error: %s", res.Error)
			}
		})
	}
}

func TestNoneOfValidator_Description(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in       []string
		expected string
	}

	testCases := map[string]testCase{
		"quoted-once": {
			in:       []string{"foo", "bar", "baz"},
			expected: `value must be none of: ["foo" "bar" "baz"]`,
		},
	}

	for name, test := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := stringvalidator.NoneOf(test.in...)

			got := v.MarkdownDescription(context.Background())

			if diff := cmp.Diff(got, test.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
