// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package listvalidator_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUniqueValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		list                types.List
		expectedDiagnostics diag.Diagnostics
	}{
		"null-list": {
			list:                types.ListNull(types.StringType),
			expectedDiagnostics: nil,
		},
		"unknown-list": {
			list:                types.ListUnknown(types.StringType),
			expectedDiagnostics: nil,
		},
		"null-value": {
			list: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringNull()},
			),
			expectedDiagnostics: nil,
		},
		"null-values-duplicate": {
			list: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringNull(), types.StringNull()},
			),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Duplicate List Value",
					"This attribute contains duplicate values of: <null>",
				),
			},
		},
		"null-values-valid": {
			list: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringNull(), types.StringValue("test")},
			),
			expectedDiagnostics: nil,
		},
		"unknown-value": {
			list: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringUnknown()},
			),
			expectedDiagnostics: nil,
		},
		"unknown-values-duplicate": {
			list: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringUnknown(), types.StringUnknown()},
			),
			expectedDiagnostics: nil,
		},
		"unknown-values-valid": {
			list: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringUnknown(), types.StringValue("test")},
			),
			expectedDiagnostics: nil,
		},
		"known-value": {
			list: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("test")},
			),
			expectedDiagnostics: nil,
		},
		"known-values-duplicate": {
			list: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("test"), types.StringValue("test")},
			),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Duplicate List Value",
					"This attribute contains duplicate values of: \"test\"",
				),
			},
		},
		"known-values-valid": {
			list: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("test1"), types.StringValue("test2")},
			),
			expectedDiagnostics: nil,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.ListRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.list,
			}
			response := validator.ListResponse{}
			listvalidator.UniqueValues().ValidateList(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
