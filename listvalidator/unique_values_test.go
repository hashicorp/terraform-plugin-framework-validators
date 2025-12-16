// Copyright IBM Corp. 2022, 2025
// SPDX-License-Identifier: MPL-2.0

package listvalidator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
)

func TestUniqueValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		list                types.List
		expectedDiagnostics diag.Diagnostics
		expectedFuncError   *function.FuncError
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
			expectedFuncError: function.NewArgumentFuncError(
				0,
				"Duplicate List Value: This attribute contains duplicate values of: <null>",
			),
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
			expectedFuncError: function.NewArgumentFuncError(
				0,
				"Duplicate List Value: This attribute contains duplicate values of: \"test\"",
			),
		},
		"multiple-known-values-duplicate": {
			list: types.ListValueMust(
				types.StringType,
				[]attr.Value{
					types.StringValue("test-val-1"),
					types.StringValue("test-val-1"),
					types.StringValue("test-val-2"),
					types.StringValue("test-val-2"),
				},
			),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Duplicate List Value",
					"This attribute contains duplicate values of: \"test-val-1\"",
				),
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Duplicate List Value",
					"This attribute contains duplicate values of: \"test-val-2\"",
				),
			},
			expectedFuncError: function.NewArgumentFuncError(
				0,
				"Duplicate List Value: This attribute contains duplicate values of: \"test-val-1\"\n"+
					"Duplicate List Value: This attribute contains duplicate values of: \"test-val-2\"",
			),
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

		t.Run(fmt.Sprintf("ValidateList - %s", name), func(t *testing.T) {
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

		t.Run(fmt.Sprintf("ValidateParameterList - %s", name), func(t *testing.T) {
			t.Parallel()

			request := function.ListParameterValidatorRequest{
				ArgumentPosition: 0,
				Value:            testCase.list,
			}
			response := function.ListParameterValidatorResponse{}
			listvalidator.UniqueValues().ValidateParameterList(context.Background(), request, &response)

			if diff := cmp.Diff(response.Error, testCase.expectedFuncError); diff != "" {
				t.Errorf("unexpected function error difference: %s", diff)
			}
		})
	}
}
